package types

type MeasureUnit string

const (
	MeasureUnitGrams       = "g"
	MeasureUnitMillilitres = "mL"
	MeasureUnitTeaspoon    = "tsp"
	MeasureUnitTablespoon  = "tbsp"
	MeasureUnitPiece       = "pc"
	MeasureUnitCup         = "cup"
)

var MeasureUnits = []MeasureUnit{
	MeasureUnitGrams,
	MeasureUnitCup,
	MeasureUnitMillilitres,
	MeasureUnitPiece,
	MeasureUnitTablespoon,
	MeasureUnitTeaspoon,
}
