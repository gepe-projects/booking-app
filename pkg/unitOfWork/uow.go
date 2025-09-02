package uow

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Unit of Work (UoW) itu sebenarnya abstraction di atas transaction handling biar service layer nggak terlalu “kotor” dengan begin/commit/rollback, tapi tetap punya kontrol penuh kalau butuh cross-domain.

// Flow UnitOfWork
// 1. Service layer panggil uow.Do(ctx, func(tx *sqlx.Tx) error { ... }).
// 2. UnitOfWork otomatis bikin BEGIN, kasih *sqlx.Tx ke callback.
// 3. Callback jalanin business logic pakai repository (repo wajib ada param *sqlx.Tx, dan bukan pake db.xx tapi tx.xx).
// 4. Kalau error → ROLLBACK. Kalau sukses → COMMIT.

type UnitOfWork interface {
	Do(ctx context.Context, fn func(tx *sqlx.Tx) error) error
}

type unitOfWork struct {
	db *sqlx.DB
}

func NewUnitOfWork(db *sqlx.DB) UnitOfWork {
	return &unitOfWork{db: db}
}

// Do menjalankan fungsi yang diberikan dalam konteks transaksi database.
// Fungsi ini akan memulai transaksi baru, memberikan objek transaksi ke fungsi yang diberikan,
// dan melakukan commit jika fungsi tidak mengembalikan error. Jika fungsi mengembalikan error
// atau gagal memulai transaksi, maka transaksi akan di-rollback dan error dikembalikan.
// Metode ini memastikan semua operasi di dalam fungsi dijalankan secara atomik.
//
// Parameter:
//
//	ctx - Context untuk mengontrol pembatalan dan batas waktu.
//	fn  - Fungsi yang menerima objek transaksi dan menjalankan operasi database.
//
// Return:
//
//	Error jika gagal memulai transaksi, jika fungsi mengembalikan error,
//	atau jika commit transaksi gagal; jika tidak, nil.
func (u *unitOfWork) Do(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
