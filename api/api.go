package api

import (
	"context"
	"net/http"
	"solana_rpc/config"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gin-gonic/gin"
)

// Send the test Token to user
func TokenAccountBalance(c *gin.Context) {
	client := rpc.New(config.SolanaRpc)
	pubKey, err := solana.PublicKeyFromBase58(c.Query("account"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": ACCOUNT_ERROR,
			"err_msg":  "account error. " + c.Query("account"),
		})
		return
	}

	var data token.Account
	spl := c.Query("spl")
	if spl == "0" {
		var out *rpc.GetBalanceResult
		if out, err = client.GetBalance(context.TODO(), pubKey, "finalized"); err == nil {
			data.Amount = out.Value
			data.Owner = pubKey
		}
	} else {
		err = client.GetAccountDataInto(context.TODO(), pubKey, &data)
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": SYSTEM_ERROR,
			"err_msg":  "rpc error. " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    true,
		"owner":     data.Owner.String(),
		"programid": data.Mint.String(),
		"amount":    data.Amount,
	})
}
