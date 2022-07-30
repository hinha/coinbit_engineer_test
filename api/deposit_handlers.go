package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"

	coin "github.com/hinha/coinbit_test"
	"github.com/hinha/coinbit_test/pb"
)

type depositRequest struct {
	WalletId string  `json:"wallet_id"`
	Amount   float64 `json:"amount"`
}

func newDepositHandlerFunc(emitter *goka.Emitter, stream goka.Stream) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req depositRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			NewErrorf(w, err, http.StatusBadRequest).ResponseJson()
			return
		}

		if err := json.Unmarshal(body, &req); err != nil {
			NewErrorf(w, err, http.StatusBadRequest).ResponseJson()
			return
		}

		if !(req.Amount > 0) {
			NewErrorf(w, errors.New("amount must greater than 0"), http.StatusBadRequest).ResponseJson()
			return
		}

		deposit := &pb.Deposit{
			WalletId: req.WalletId,
			Amount:   req.Amount,
		}

		if stream == coin.DepositStream {
			err = emitter.EmitSync(req.WalletId, deposit)
		} else {
			deposit.Amount = -1 * deposit.Amount
			err = emitter.EmitSync(req.WalletId, deposit)
		}
		if err != nil {
			NewErrorf(w, err, http.StatusBadRequest).ResponseJson()
			return
		}

		ResponseJson(w, "success", nil, http.StatusOK)
	}
}

type depositResponse struct {
	WalletId       string  `json:"wallet_id"`
	Balance        float64 `json:"balance"`
	AboveThreshold bool    `json:"above_threshold"`
}

func newCheckHandlerFunc(emitter *goka.Emitter, balanceView, flagView *goka.View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			totalBalance   float64
			aboveThreshold bool
		)

		walletID := mux.Vars(r)["wallet_id"]
		response := depositResponse{
			WalletId:       walletID,
			Balance:        totalBalance,
			AboveThreshold: aboveThreshold,
		}

		val, _ := balanceView.Get(walletID)
		if val == nil {
			ResponseJson(w, "success", response, http.StatusOK)
			return
		}

		// if v, ok := val.(*pb.DepositHistory); ok {
		// 	for _, m := range v.Deposits {
		// 		totalBalance += m.Amount
		// 	}
		// }

		dp := val.(*pb.DepositHistory)
		for _, m := range dp.Deposits {
			totalBalance += m.Amount
		}

		flaggerVal, _ := flagView.Get(walletID)
		if flaggerVal != nil {
			b := flaggerVal.(*pb.FlagValue)
			aboveThreshold = b.Flagged
		}

		response.Balance = totalBalance
		response.AboveThreshold = aboveThreshold
		ResponseJson(w, "success", response, http.StatusOK)
	}
}
