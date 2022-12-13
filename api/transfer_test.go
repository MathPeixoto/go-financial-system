package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/MathPeixoto/go-financial-system/db/mock"
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/MathPeixoto/go-financial-system/token"
	"github.com/MathPeixoto/go-financial-system/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateTransferAPI(t *testing.T) {
	// prepare test data
	userOne, _ := randomUser(t)
	userTwo, _ := randomUser(t)
	userThree, _ := randomUser(t)

	accountOne := brlAccount(userOne.Username)
	accountTwo := brlAccount(userTwo.Username)
	accountThree := usdAccount(userThree.Username)

	validTransferRequest := transferRequest{
		FromAccountID: accountOne.ID,
		ToAccountID:   accountTwo.ID,
		Amount:        1000,
		Currency:      util.BRL,
	}

	invalidTransferWithDifferentCurrenciesRequest := transferRequest{
		FromAccountID: accountOne.ID,
		ToAccountID:   accountThree.ID,
		Amount:        1000,
		Currency:      accountThree.Currency,
	}

	invalidTransferRequest := transferRequest{
		FromAccountID: accountOne.ID,
		ToAccountID:   accountTwo.ID,
		Amount:        1000,
		Currency:      "ABC",
	}

	validTransferTxParams := getTransferParams(validTransferRequest)

	dbTransfer := createTransferTx(validTransferTxParams)

	testCases := []struct {
		name          string
		body          any
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: validTransferRequest,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.FromAccountID)).Times(1).Return(accountOne, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.ToAccountID)).Times(1).Return(accountTwo, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(validTransferTxParams)).Times(1).Return(dbTransfer, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransferTx(t, recorder.Body, dbTransfer)
			},
		},
		{
			name: "BadRequest",
			body: invalidTransferRequest,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest - Invalid Account One",
			body: validTransferRequest,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.FromAccountID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest - Invalid Account Two",
			body: validTransferRequest,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.FromAccountID)).Times(1).Return(accountOne, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.ToAccountID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest - Invalid Currency Between Accounts",
			body: invalidTransferWithDifferentCurrenciesRequest,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(invalidTransferWithDifferentCurrenciesRequest.FromAccountID)).Times(1).Return(accountOne, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error - Could not validate account",
			body: validTransferRequest,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.FromAccountID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal Server Error - Could not create transfer",
			body: validTransferRequest,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.FromAccountID)).Times(1).Return(accountOne, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.ToAccountID)).Times(1).Return(accountTwo, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(validTransferTxParams)).Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized - different user logged in",
			body: validTransferRequest,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userTwo.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.FromAccountID)).Times(1).Return(accountOne, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			testBody := testCase.body
			body, err := json.Marshal(testBody)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(body))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestGetTransfer(t *testing.T) {
	// prepare test data
	userOne, _ := randomUser(t)
	userTwo, _ := randomUser(t)

	accountOne := randomAccount(userOne.Username)
	accountTwo := randomAccount(userTwo.Username)

	dbTransfer := randomTransfer(accountOne, accountTwo)

	testCases := []struct {
		name          string
		id            int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   dbTransfer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTransfer(gomock.Any(), gomock.Eq(dbTransfer.ID)).Times(1).Return(dbTransfer, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransfer(t, recorder.Body, dbTransfer)
			},
		},
		{
			name: "NotFound",
			id:   dbTransfer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTransfer(gomock.Any(), gomock.Eq(dbTransfer.ID)).Times(1).Return(db.Transfer{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			id:   -1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			id:   dbTransfer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, request, tokenMaker, authTypeBearer, userOne.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTransfer(gomock.Any(), gomock.Eq(dbTransfer.ID)).Times(1).Return(db.Transfer{}, sql.ErrConnDone)
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/transfers/%d", testCase.id), nil)

			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func requireBodyMatchTransferTx(t *testing.T, body *bytes.Buffer, transfer db.TransferTxResult) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTransferTx db.TransferTxResult
	err = json.Unmarshal(data, &gotTransferTx)
	require.NoError(t, err)

	require.Equal(t, transfer, gotTransferTx)
}

func requireBodyMatchTransfer(t *testing.T, body *bytes.Buffer, transfer db.Transfer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTransfer db.Transfer
	err = json.Unmarshal(data, &gotTransfer)
	require.NoError(t, err)

	require.Equal(t, transfer, gotTransfer)
}
