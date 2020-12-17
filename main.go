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

func owns(kittyID uint64, owner string) (bool, error) {
	return false, nil
}

func approvedFor(kittyID uint64, account string) (bool, error) {
	return false, nil
}

func approve(kittdyID uint64, account string) error {
	return nil
}

func (c *KittyContract) Transfer(ctx contractapi.TransactionContextInterface, to string, kittyID uint64) error {
	return nil
}

func (c *KittyContract) Approve(ctx contractapi.TransactionContextInterface, kittdyID uint64, account string) error {
	return nil
}

func (c *KittyContract) TotalSuppy(ctx contractapi.TransactionContextInterface) uint64 {
	return 0
}

func (c *KittyContract) OwnerOf(ctx contractapi.TransactionContextInterface, kittyID uint64) (string, error) {
	return "", nil
}

func (c *KittyContract) TokensOfOwner(owner string) ([]uint64, error) {
	return []uint64{}, nil
}

func (c *KittyContract) PregnantKitties() (uint64, error) {
	return 0, nil
}

func isReadyToGiveBirth(matron Kitty) (bool, error) {
	return true, nil
}

func isReadyToBreed(kittyID uint64) (bool, error) {
	return false, nil
}

func (c *KittyContract) IsReadyToBreed(kittyID uint64) (bool, error) {
	return false, nil
}

func isSiringPermitted(matronID, sireID uint64) (bool, error) {
	return true, nil
}

func triggerCooldown(kittyID uint64) error {
	return nil
}

func (c *KittyContract) ApproveSiring(ctx contractapi.TransactionContextInterface, kittyID uint64, siringPartner string) error {
	return nil
}

func isPregnant(kittyID uint64) (bool, error) {
	return false, nil
}

func isValidMatingPair(sire uint64, matron uint64) (bool, error) {
	return true, nil
}

func (c *KittyContract) CanBreedWith(ctx contractapi.TransactionContextInterface, sire uint64, matron uint64) (bool, error) {
	return true, nil
}

func breedWith(sire uint64, matron uint64) error {
	return nil
}

func (c *KittyContract) BreedWithAuto(ctx contractapi.TransactionContextInterface, sireID uint64, matronID uint64) error {
	return nil
}

func (c *KittyContract) GiveBirth(ctx contractapi.TransactionContextInterface, matronID uint64) error {
	return nil
}

func (c *KittyContract) TransferFrom(ctx contractapi.TransactionContextInterface, from, to string, kittyID uint64) error {
	// TODO (hb237): extend the functionality of this function to match the ethereum implementation
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
