package http

import (
	"errors"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/types"
)

func (srv *Server) queue_proposal(c echo.Context) error {
	// Unmarshal request body
	reqBody := new(types.QueueProposalRequestBody)
	err := c.Bind(reqBody)
	if err != nil {
		srv.returnError(c, http.StatusUnprocessableEntity, err)
	}

	// Validate proposal request
	if reqBody.TxDest == (common.Address{}) {
		return srv.returnError(c, http.StatusBadRequest, errors.New("require non zero transaction destination"))
	}
	if reqBody.TxData == nil || len(reqBody.TxData) == 0 {
		return srv.returnError(c, http.StatusBadRequest, errors.New("require non empty transaction data"))
	}
	if reqBody.TxList == nil || len(reqBody.TxList) == 0 {
		return srv.returnError(c, http.StatusBadRequest, errors.New("require non empty transaction list"))
	}

	// Publish to rabbitmq
	err = srv.queue.Publish(c.Request().Context(), *reqBody)
	if err != nil {
		return srv.returnError(c, http.StatusInternalServerError, errors.New("unable to queue proposal"))
	}

	return c.JSON(http.StatusOK, reqBody)
}
