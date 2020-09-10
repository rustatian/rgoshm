package semaphore

import (
	"errors"
	"syscall"
	"unsafe"
)

/* The following System V style IPC functions implement a semaphore
   handling.  The definition is found in XPG2.  */

/* Structure used for argument to `semop' to describe operations.  */
//struct sembuf
//{
//unsigned short int sem_num;	/* semaphore number */
//short int sem_op;		/* semaphore operation */
//short int sem_flg;		/* operation flag */
//};

//The sem_num field identifies the semaphore within the set upon which the opera-
//tion is to be performed. The sem_op field specifies the operation to be performed:
//
//If sem_op is greater than 0, the value of sem_op is added to the semaphore value.
//As a result, other processes waiting to decrease the semaphore value may be
//awakened and perform their operations. The calling process must have alter
//(write) permission on the semaphore.

//If sem_op equals 0, the value of the semaphore is checked to see whether it cur-
//rently equals 0. If it does, the operation completes immediately; otherwise,
//semop() blocks until the semaphore value becomes 0. The calling process must
//have read permission on the semaphore.

//If sem_op is less than 0, decrease the value of the semaphore by the amount
//specified in sem_op. If the current value of the semaphore is greater than or
//equal to the absolute value of sem_op, the operation completes immediately.
//Otherwise, semop() blocks until the semaphore value has been increased to a
//level that permits the operation to be performed without resulting in a nega-
//tive value. The calling process must have alter permission on the semaphore.

type Flag int

const (
	IPC_CREAT  Flag = 01000 /* Create key if key does not exist. */
	IPC_EXCL   Flag = 02000 /* Fail if key exists.  */
	IPC_NOWAIT Flag = 04000

	IPC_RMID Flag = 0 /* Remove identifier.  */
	IPC_SET  Flag = 1 /* Set `ipc_perm' options.  */
	IPC_STAT Flag = 2 /* Get `ipc_perm' options.  */
)

type Semaphore struct {
	id            int
	numSemaphores int
}

type sembuf struct {
	sem_num uint16 // according to the standard, unsigned short should be at least 2 bytes
	sem_op  int16  // /* Operation to be performed */
	sem_flg int16  // /* Operation flags (IPC_NOWAIT and SEM_UNDO) */
}

func NewSemaphore() *Semaphore {
	return &Semaphore{}
}

// Wait waits until data will be written (semaphore unlocked in the other process)
func (s *Semaphore) Wait() error {
	return nil
}

// int semget(key_t key , int nsems , int semflg );
// flags which can be used:
// IPC_CREAT - If no semaphore set with the specified key exists, create a new set.
// IPC_EXCL If IPC_CREAT was also specified, and a semaphore set with the specified key already exists, fail with the error EEXIST
// return semaphore ID
func (s *Semaphore) GetValue(key string, numSems int, flags ...Flag) (int, error) {
	path, err := syscall.BytePtrFromString(key)
	if err != nil {
		return -1, err
	}

	// OR flags
	var flgs Flag
	for i := 0; i < len(flags); i++ {
		flgs = flgs | flags[i]
	}
	flgs = flgs | 0666

	val, _, errno := syscall.Syscall(syscall.SYS_SEMGET, uintptr(unsafe.Pointer(path)), uintptr(numSems), uintptr(flgs))
	if errno != 0 {
		return -1, errors.New(errno.Error())
	}
	//
	s.id = int(val)

	return int(val), nil
}