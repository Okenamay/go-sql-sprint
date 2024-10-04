package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	// Добавление строки в таблицу parcel с использованием данных из переменной p:
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	// Определение идентификатора последней добавленной строки:
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Вывод идентификатора последней добавленной строки:
	return int(id), nil // LastInsertId() выдаёт int64, пришлось конвертировать:
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	// Чтение одной строки из таблицы parcel по заданному number:
	p := Parcel{}
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	// Отладка:
	if err != nil {
		// Отладка:
		fmt.Println("Ошибка в функции Get:", err)
	}

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	// Чтение строк из таблицы parcel по заданному client
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Заполнение среза []Parcel данными из таблицы:
	var res []Parcel

	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println("Ошибка в функции GetByClient:", err)
			return nil, err
		}

		res = append(res, p)
	}

	return res, err
}

func (s ParcelStore) SetStatus(number int, status string) error {

	// Обновление статуса в таблице parcel:
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {

	// Обновление адреса в таблице parcel:
	// Проверим, имеет ли статус записи с номером number значение registered:
	p := Parcel{}
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number",
		sql.Named("number", number))
	err := row.Scan(&p.Status)
	if err != nil {
		// Отладка:
		fmt.Println("Ошибка в функции SetAddress:", err)
		return err
	}

	// Если условние соблюдено, обновим адрес в таблице parcel:
	if p.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
		return err
	}

	return err
}

func (s ParcelStore) Delete(number int) error {
	// Удаление строки из таблицы parcel:
	// Проверим, имеет ли статус записи с номером number значение registered:
	p := Parcel{}
	rows := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number",
		sql.Named("number", number))
	err := rows.Scan(&p.Status)
	if err != nil {
		// Отладка:
		fmt.Println("Ошибка в функции Delete:", err)
		return err
	}

	// Если условние соблюдено, удалим строку в таблице parcel:
	if p.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number",
			sql.Named("number", number))
		return err
	}

	return nil
}
