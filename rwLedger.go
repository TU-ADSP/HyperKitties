package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func readFromLedger(ctx contractapi.TransactionContextInterface, key string) error {
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}

	if key == kittyIndexToOwnerNAME {
		err = json.Unmarshal(assetJSON, &kittyIndexToOwner)
	}
	if key == kittyIndexToApprovedNAME {
		err = json.Unmarshal(assetJSON, &kittyIndexToApproved)
	}
	if key == sireAllowedToAddressNAME {
		err = json.Unmarshal(assetJSON, &sireAllowedToAddress)
	}
	if key == kittiesNAME {
		err = json.Unmarshal(assetJSON, &kitties)
	}
	if err != nil {
		return err
	}

	return nil
}

func writeToLedger(ctx contractapi.TransactionContextInterface, key string) error {
	var assetJSON []byte
	var err error

	if key == kittyIndexToOwnerNAME {
		assetJSON, err = json.Marshal(kittyIndexToOwner)
	}
	if key == kittyIndexToApprovedNAME {
		assetJSON, err = json.Marshal(kittyIndexToApproved)
	}
	if key == sireAllowedToAddressNAME {
		assetJSON, err = json.Marshal(sireAllowedToAddress)
	}
	if key == kittiesNAME {
		assetJSON, err = json.Marshal(kitties)
	}
	if err != nil {
		return err
	}

	if err = ctx.GetStub().PutState(key, assetJSON); err != nil {
		return err
	}

	return nil
}
