package main

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
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

var kitties = KittyList{Kitty{}}

const kittiesNAME = "kitties"

var kittyIndexToOwner = []string{""}

const kittyIndexToOwnerNAME = "kittyIndexToOwner"

var kittyIndexToApproved = []string{""}

const kittyIndexToApprovedNAME = "kittyIndexToApproved"

var sireAllowedToAddress = []string{""}

const sireAllowedToAddressNAME = "kittyIndexToAddress"

var g_event = map[string]interface{}{}

func transfer(ctx contractapi.TransactionContextInterface, from, to string, kittyID uint64) error {
	kittyIndexToOwner[kittyID] = to

	kittyIndexToApproved[kittyID] = ""
	sireAllowedToAddress[kittyID] = ""

	payload := map[string]interface{}{"from": from, "to": to, "kittyID": kittyID}
	g_event["Transfer"] = payload

	return nil
}

func createKitty(ctx contractapi.TransactionContextInterface, matronID, sireID, generation, genes uint64, owner string) (uint64, error) {
	cooldownIndex := uint8(generation / 2)
	if cooldownIndex > 13 {
		cooldownIndex = 13
	}

	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return 0, err
	}
	now, err := ptypes.Timestamp(txTimestamp)
	if err != nil {
		return 0, err
	}

	kitty := Kitty{
		Genes:         genes,
		BirthTime:     now,
		CooldownEnd:   now.Add(cooldowns[cooldownIndex]),
		MatronID:      matronID,
		SireID:        sireID,
		SiringWithID:  0,
		CooldownIndex: cooldownIndex,
		Generation:    generation,
	}
	kitties = append(kitties, kitty)
	newKittenID := uint64(len(kitties) - 1)

	kittyIndexToOwner = append(kittyIndexToOwner, "")
	kittyIndexToApproved = append(kittyIndexToApproved, "")
	sireAllowedToAddress = append(sireAllowedToAddress, "")

	payload := map[string]interface{}{"owner": owner, "newKittenID": newKittenID, "matronID": matronID, "sireID": sireID, "genes": genes}
	g_event["Birth"] = payload

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

func (c *KittyContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("Entering InitLedger: Doing essentially nothing...")
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ktct := KittyContract{}
	ktct.BeforeTransaction = BeforeTransaction
	ktct.AfterTransaction = AfterTransaction
	assetChaincode, err := contractapi.NewChaincode(&ktct)
	if err != nil {
		log.Panicf("Error creating HyperKitty chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting HyperKitty chaincode: %v", err)
	}
}
