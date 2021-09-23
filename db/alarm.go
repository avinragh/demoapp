package db

import uuid "github.com/satori/go.uuid"

func (db *DB) FindAlarmById(id string) (*Alarm, error) {
	txn := db.MemDB.Txn(false)
	defer txn.Abort()

	item, err := txn.First("alarm", "id", id)
	if err != nil {
		return nil, err
	}
	alarm := item.(*Alarm)

	return alarm, nil
}

func (db *DB) FindAlarms(alarmType, resourceId *string) ([]*Alarm, error) {

	alarms := []*Alarm{}
	txn := db.MemDB.Txn(false)
	defer txn.Abort()

	if alarmType != nil && resourceId != nil {
		it, err := txn.Get("alarm", "id", "alarmType", alarmType, "resourceId", resourceId)
		if err != nil {
			return nil, err
		}
		for obj := it.Next(); obj != nil; obj = it.Next() {
			alarm := obj.(*Alarm)
			alarms = append(alarms, alarm)

		}

	} else {
		it, err := txn.Get("alarm", "id")
		if err != nil {
			return nil, err
		}
		for obj := it.Next(); obj != nil; obj = it.Next() {
			alarm := obj.(*Alarm)
			alarms = append(alarms, alarm)
		}
	}
	return alarms, nil

}

func (db *DB) AddAlarms(alarms []*Alarm) ([]*Alarm, error) {
	for _, alarm := range alarms {
		db.AddAlarm(alarm)
	}
	return alarms, nil
}

func (db *DB) AddAlarm(alarm *Alarm) (*Alarm, error) {
	if alarm.Id == nil {
		uuid := uuid.NewV4().String()
		alarm.Id = &uuid
	}
	txn := db.MemDB.Txn(true)

	if err := txn.Insert("alarm", alarm); err != nil {
		return nil, err
	}
	txn.Commit()
	return alarm, nil
}

func (db *DB) DeleteAlarm(id string) (*Alarm, error) {
	alarm, err := db.FindAlarmById(id)
	if err != nil {
		return nil, err
	}
	txn := db.MemDB.Txn(true)
	txn.Delete("alarm", alarm)
	txn.Commit()
	return alarm, nil
}
