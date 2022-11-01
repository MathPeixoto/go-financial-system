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

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
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
		t.Run(testCase.name, func(t *testing.T) {
			url := fmt.Sprintf("/accounts/%d", testCase.accountID)
			// prepare stubs
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			// start test server
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	valIdAccountRequest := CreateAccountRequest{
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
	}

	invalidIdAccountRequest := CreateAccountRequest{
		Owner:    util.RandomOwner(),
		Currency: "invalid",
	}

	accountParams := getAccountParams(valIdAccountRequest)

	dbAccount := createAccount(accountParams)

	testCases := []struct {
		name          string
		body          any
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: valIdAccountRequest,
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
			body:       invalidIdAccountRequest,
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: valIdAccountRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(accountParams)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// prepare stubs
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			// start test server
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			testBody := testCase.body
			body, err := json.Marshal(testBody)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestListAccountAPI(t *testing.T) {
	validListAccountsRequest := ListAccountsRequest{
		Limit:  5,
		Offset: 5,
	}

	invalidListAccountsRequest := ListAccountsRequest{
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

	testCases := []struct {
		name                string
		body                any
		ListAccountsRequest ListAccountsRequest
		buildStubs          func(store *mockdb.MockStore)
		checkResponse       func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:                "OK",
			ListAccountsRequest: validListAccountsRequest,
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
			ListAccountsRequest: invalidListAccountsRequest,
			buildStubs:          func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:                "InternalError",
			ListAccountsRequest: validListAccountsRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			url := fmt.Sprintf("/accounts?limit=%d&offset=%d", testCase.ListAccountsRequest.Limit, testCase.ListAccountsRequest.Offset)
			// prepare stubs
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			// start test server
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestUpdateAccountAPI(t *testing.T) {
	dbAccount := randomAccount()

	validIdAccountRequest := IdAccountRequest{
		ID: dbAccount.ID,
	}

	invalidIdAccountRequest := IdAccountRequest{
		ID: -1,
	}

	valIdAccountRequest := UpdateAccountBalanceRequest{
		Amount: util.RandomMoney(),
	}

	invalidAccountRequest := UpdateAccountBalanceRequest{
		Amount: -1,
	}

	arg := db.AddAccountBalanceParams{
		ID:     dbAccount.ID,
		Amount: valIdAccountRequest.Amount,
	}

	updatedAccount := updatedAccount(dbAccount, arg)

	testCases := []struct {
		name          string
		id            int64
		body          any
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   validIdAccountRequest.ID,
			body: valIdAccountRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddAccountBalance(gomock.Any(), gomock.Eq(arg)).Times(1).Return(updatedAccount, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, updatedAccount)
			},
		},
		{
			name:       "BadRequest",
			id:         invalidIdAccountRequest.ID,
			body:       valIdAccountRequest,
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			id:         validIdAccountRequest.ID,
			body:       invalidAccountRequest,
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   validIdAccountRequest.ID,
			body: valIdAccountRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().AddAccountBalance(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// prepare stubs
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			// start test server
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			testBody := testCase.body
			body, err := json.Marshal(testBody)
			require.NoError(t, err)

			url := fmt.Sprintf("/accounts/%d", testCase.id)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestDeleteAccountAPI(t *testing.T) {
	ID := util.RandomInt(1, 1000)

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
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
		t.Run(testCase.name, func(t *testing.T) {
			url := fmt.Sprintf("/accounts/%d", testCase.accountID)
			// prepare stubs
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			// start test server
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodDelete, url, nil)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
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
