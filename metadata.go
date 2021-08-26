package dbbin

type Metadata struct {
	Size       uint16 // tamanho do metadados
	FieldsSize uint16 // tamanho da string de referencia dos campos
	RecordSize uint16 // em bytes
	Created    uint64 // timestamp ms de criação
}
