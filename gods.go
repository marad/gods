package gods

type Value interface {
}

type Seq interface {
	Next() (Value, error)
}
