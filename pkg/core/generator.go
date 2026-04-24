package core

import (
	mydriver "gendata/pkg/driver"
)

func GenerateDefaultTableData(count int) []*mydriver.DefaultTableRow {
	rows := make([]*mydriver.DefaultTableRow, count)
	for i := 0; i < count; i++ {
		generator := newDefaultTableGenerator()
		rows[i] = generator.generateDefaultTableRow()
	}
	return rows
}

type defaultTableGenerator struct {
	fieldGen *DefaultGenerator
}

func newDefaultTableGenerator() *defaultTableGenerator {
	return &defaultTableGenerator{
		fieldGen: NewDefaultGenerator(),
	}
}

func (dtg *defaultTableGenerator) generateDefaultTableRow() *mydriver.DefaultTableRow {
	return &mydriver.DefaultTableRow{
		UserID:       dtg.fieldGen.GenerateUserID(),
		Name:         dtg.fieldGen.GenerateName(),
		Phone:        dtg.fieldGen.GeneratePhone(),
		Gender:       dtg.fieldGen.GenerateGender(),
		Age:          dtg.fieldGen.GenerateAge(),
		Birthday:     dtg.fieldGen.GenerateBirthday(),
		Email:        dtg.fieldGen.GenerateEmail(),
		State:        dtg.fieldGen.GenerateState(),
		ZipCode:      dtg.fieldGen.GenerateZipCode(),
		Nationality:  dtg.fieldGen.GenerateNationality(),
		Height:       dtg.fieldGen.GenerateHeight(),
		Weight:       dtg.fieldGen.GenerateWeight(),
		BloodType:    dtg.fieldGen.GenerateBloodType(),
		Account:      dtg.fieldGen.GenerateAccount(),
		AccountName:  dtg.fieldGen.GenerateAccountName(),
		Password:     dtg.fieldGen.GeneratePassword(),
		OnlineStatus: dtg.fieldGen.GenerateOnlineStatus(),
		CreatedAt:    dtg.fieldGen.GenerateDateTime(),
	}
}
