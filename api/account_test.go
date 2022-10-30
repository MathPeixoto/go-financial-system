package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/MathPeixoto/go-financial-system/db/mock"
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/MathPeixoto/go-financial-system/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testScenarios struct {
	name                string
	accountID           int64
	checkResponse       func(recorder *httptest.ResponseRecorder)
	buildStubs          func(store *mockdb.MockStore)
	body                any
	listAccountsRequest listAccountsRequest
}

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []testScenarios{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:       "BadRequest",
			accountID:  -1,
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		url := fmt.Sprintf("/accounts/%d", testCase.accountID)
		runTestCases(t, testCase, http.MethodGet, url)
	}
}

func TestCreateAccountAPI(t *testing.T) {
	validAccountRequest := createAccountRequest{
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
	}

	invalidAccountRequest := createAccountRequest{
		Owner:    util.RandomOwner(),
		Currency: "invalid",
	}

	accountParams := getAccountParams(validAccountRequest)

	dbAccount := account(accountParams)

	testCases := []testScenarios{
		{
			name: "OK",
			body: validAccountRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(accountParams)).Times(1).Return(dbAccount, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, dbAccount)
			},
		},
		{
			name:       "BadRequest",
			body:       invalidAccountRequest,
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: validAccountRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(accountParams)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		runTestCases(t, testCase, http.MethodPost, "/accounts")
	}
}

func TestListAccountAPI(t *testing.T) {
	validListAccountsRequest := listAccountsRequest{
		Limit:  5,
		Offset: 5,
	}

	invalidListAccountsRequest := listAccountsRequest{
		Limit:  -1,
		Offset: -1,
	}

	arg := db.ListAccountsParams{
		Limit:  validListAccountsRequest.Limit,
		Offset: (validListAccountsRequest.Offset - 1) * validListAccountsRequest.Limit,
	}

	var accounts = []db.Account{
		randomAccount(),
		randomAccount(),
		randomAccount(),
		randomAccount(),
		randomAccount(),
	}

	testCases := []testScenarios{
		{
			name:                "OK",
			listAccountsRequest: validListAccountsRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(accounts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name:                "BadRequest",
			listAccountsRequest: invalidListAccountsRequest,
			buildStubs:          func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:                "InternalError",
			listAccountsRequest: validListAccountsRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		url := fmt.Sprintf("/accounts?limit=%d&offset=%d", testCase.listAccountsRequest.Limit, testCase.listAccountsRequest.Offset)
		runTestCases(t, testCase, http.MethodGet, url)
	}
}

func TestUpdateAccountAPI(t *testing.T) {
	// TODO
}

func TestDeleteAccountAPI(t *testing.T) {
	ID := util.RandomInt(1, 1000)

	testCases := []testScenarios{
		{
			name:      "OK",
			accountID: ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(ID)).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			accountID:  -1,
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(ID)).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		url := fmt.Sprintf("/accounts/%d", testCase.accountID)
		runTestCases(t, testCase, http.MethodDelete, url)
	}
}

func runTestCases(t *testing.T, testCase testScenarios, method string, url string) {
	t.Run(testCase.name, func(t *testing.T) {
		// prepare stubs
		ctrl := gomock.NewController(t)
		store := mockdb.NewMockStore(ctrl)
		testCase.buildStubs(store)

		// start test server
		server := NewServer(store)
		recorder := httptest.NewRecorder()
		var request = &http.Request{}
		var err error

		if testBody := testCase.body; testBody != nil {
			body, err := json.Marshal(testBody)
			require.NoError(t, err)
			request, err = http.NewRequest(method, url, bytes.NewReader(body))
		} else {
			request, err = http.NewRequest(method, url, nil)
		}

		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)
		testCase.checkResponse(recorder)
	})
}

func getAccountParams(args createAccountRequest) db.CreateAccountParams {
	return db.CreateAccountParams{
		Owner:    args.Owner,
		Currency: args.Currency,
		Balance:  0,
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func account(accountParams db.CreateAccountParams) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    accountParams.Owner,
		Balance:  accountParams.Balance,
		Currency: accountParams.Currency,
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)

	require.Equal(t, accounts, gotAccounts)
}
