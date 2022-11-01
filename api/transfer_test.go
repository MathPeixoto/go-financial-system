package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

func TestCreateTransferAPI(t *testing.T) {
	accountOne := brlAccount()
	accountTwo := brlAccount()
	accountThree := usdAccount()

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
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: validTransferRequest,
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
			name:       "BadRequest",
			body:       invalidTransferRequest,
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest - Invalid Account One",
			body: validTransferRequest,
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
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.FromAccountID)).Times(1).Return(accountOne, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(validTransferRequest.ToAccountID)).Times(1).Return(accountTwo, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(validTransferTxParams)).Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)
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

			request, err := http.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func requireBodyMatchTransferTx(t *testing.T, body *bytes.Buffer, transfer db.TransferTxResult) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTransfer db.TransferTxResult
	err = json.Unmarshal(data, &gotTransfer)
	require.NoError(t, err)

	require.Equal(t, transfer, gotTransfer)
}
