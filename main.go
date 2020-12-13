package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type KittyContract struct {
	contractapi.Contract
}

type Kitty struct {
	Genes         uint64    `json:"genes"`
	BirthTime     time.Time `json:"birth_time"`
	CooldownEnd   time.Time `json:"cooldown_end"`
	MatronID      uint64    `json:"matron_id"`
	SireID        uint64    `json:"sire_id"`
	SiringWithID  uint64    `json:"siring_with_id"`
	CooldownIndex uint8     `json:"cooldown_index"`
	Generation    uint64    `json:"generation"`
}

type KittyList []Kitty

var cooldowns = []time.Duration{
	1 * time.Second,
	2 * time.Second,
	5 * time.Second,
	10 * time.Second,
	30 * time.Second,
	1 * time.Minute,
	2 * time.Minute,
	4 * time.Minute,
	8 * time.Minute,
	16 * time.Minute,
	1 * time.Hour,
	2 * time.Hour,
	4 * time.Hour,
	7 * time.Hour,
}

var kitties KittyList = KittyList{Kitty{}}

const kittiesNAME = "kitties"

var kittyIndexToOwner = []string{""}

const kittyIndexToOwnerNAME = "kittyIndexToOwner"

var kittyIndexToApproved = []string{""}

const kittyIndexToApprovedNAME = "kittyIndexToApproved"

var sireAllowedToAddress = []string{""}

const sireAllowedToAddressNAME = "kittyIndexToAddress"

func transfer(ctx contractapi.TransactionContextInterface, from, to string, kittyID uint64) error {
	if err := readFromLedger(ctx, kittyIndexToOwnerNAME); err != nil {
		return err
	}
	kittyIndexToOwner[kittyID] = to
	if err := writeToLedger(ctx, kittyIndexToOwnerNAME); err != nil {
		return err
	}

	if from != "" {
		if err := readFromLedger(ctx, kittyIndexToApprovedNAME); err != nil {
			return err
		}
		kittyIndexToApproved[kittyID] = ""
		if err := writeToLedger(ctx, kittyIndexToApprovedNAME); err != nil {
			return err
		}
		if err := readFromLedger(ctx, sireAllowedToAddressNAME); err != nil {
			return err
		}
		sireAllowedToAddress[kittyID] = ""
		if err := writeToLedger(ctx, sireAllowedToAddressNAME); err != nil {
			return err
		}
	}

	payload := map[string]interface{}{"from": from, "to": to, "kittyID": kittyID}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().SetEvent("Transfer", jsonPayload); err != nil {
		return err
	}

	return nil
}

func createKitty(ctx contractapi.TransactionContextInterface, matronID, sireID, generation, genes uint64, owner string) (uint64, error) {
	if err := readFromLedger(ctx, kittiesNAME); err != nil {
		return 0, err
	}
	cooldownIndex := uint8(generation / 2)
	if cooldownIndex > 13 {
		cooldownIndex = 13
	}

	kitty := Kitty{
		Genes:         genes,
		BirthTime:     time.Now(),
		CooldownEnd:   time.Now().Add(cooldowns[cooldownIndex]),
		MatronID:      matronID,
		SireID:        sireID,
		SiringWithID:  0,
		CooldownIndex: cooldownIndex,
		Generation:    generation,
	}
	kitties = append(kitties, kitty)
	newKittenID := uint64(len(kitties) - 1)

	payload := map[string]interface{}{"owner": owner, "newKittenID": newKittenID, "matronID": matronID, "sireID": sireID, "genes": genes}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	if err := ctx.GetStub().SetEvent("Birth", jsonPayload); err != nil {
		return 0, err
	}

	if err := transfer(ctx, "", owner, newKittenID); err != nil {
		return 0, err
	}

	return newKittenID, nil
}

func (c *KittyContract) Transfer(ctx contractapi.TransactionContextInterface, from, to string, kittyID uint64) error {
	return transfer(ctx, from, to, kittyID)
}

func (c *KittyContract) CreateKitty(ctx contractapi.TransactionContextInterface, matronID, sireID, generation, genes uint64, owner string) error {
	_, err := createKitty(ctx, matronID, sireID, generation, genes, owner)
	return err
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&KittyContract{})
	if err != nil {
		log.Panicf("Error creating HyperKitty chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting HyperKitty chaincode: %v", err)
	}
}
