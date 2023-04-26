package utils

import (
	"fmt"
	"github.com/sony/sonyflake"
	"time"
)

var (
	sonyFlake     *sonyflake.Sonyflake
	sonyMachineId uint16
)

func getMachineId() (uint16, error) {
	return sonyMachineId, nil
}

func InitSnowFlake(machineId uint16) error {
	sonyMachineId = machineId
	// init start time
	t, _ := time.Parse("2006-01-02", "2023-01-01")
	settings := sonyflake.Settings{
		StartTime: t,
		MachineID: getMachineId,
	}
	sonyFlake = sonyflake.NewSonyflake(settings)
	return nil
}

func GetSnowFlake() (id uint64, err error) {
	if sonyFlake == nil {
		err = fmt.Errorf("SnowFlake Init Error: nil")
		return
    }
	id, err = sonyFlake.NextID()
	return
}
