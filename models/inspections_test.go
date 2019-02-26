package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CheckCoherenceSuiteConstats_ErrSuiteManquante(t *testing.T) {
	assert := require.New(t)
	inspection := Inspection{
		PointsDeControle: []PointDeControle{
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
		},
	}
	assert.Equal(ErrSuiteManquante, inspection.CheckCoherenceSuiteConstats())
}

func Test_CheckCoherenceSuiteConstats_ErrConstatManquant(t *testing.T) {
	assert := require.New(t)
	inspection := Inspection{
		PointsDeControle: []PointDeControle{
			PointDeControle{
				Publie: true,
			},
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
		},
		Suite: &Suite{
			Type: TypeSuiteAucune,
		},
	}
	assert.Equal(ErrConstatManquant, inspection.CheckCoherenceSuiteConstats())
}

func Test_CheckCoherenceSuiteConstats_ErrPointDeControleNonPublie(t *testing.T) {
	assert := require.New(t)
	inspection := Inspection{
		PointsDeControle: []PointDeControle{
			PointDeControle{
				Publie: false,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
		},
		Suite: &Suite{
			Type: TypeSuiteAucune,
		},
	}
	assert.Equal(ErrPointDeControleNonPublie, inspection.CheckCoherenceSuiteConstats())
}

func TestCheckCoherenceSuiteConformeEtPointsDeControleConformes(t *testing.T) {
	assert := require.New(t)
	inspection := Inspection{
		PointsDeControle: []PointDeControle{
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
		},
		Suite: &Suite{
			Type: TypeSuiteAucune,
		},
	}
	assert.Equal(nil, inspection.CheckCoherenceSuiteConstats())
}

func TestCheckCoherenceSuiteNonConformeEtPointDeControleNonConforme(t *testing.T) {
	assert := require.New(t)
	inspection := Inspection{
		PointsDeControle: []PointDeControle{
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatObservation,
				},
			},
		},
		Suite: &Suite{
			Type: TypeSuiteObservation,
		},
	}
	assert.Equal(nil, inspection.CheckCoherenceSuiteConstats())
}

func TestCheckIncoherenceSuiteNonConformeEtPointsDeControleConformes(t *testing.T) {
	assert := require.New(t)
	inspection := Inspection{
		PointsDeControle: []PointDeControle{
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
		},
		Suite: &Suite{
			Type: TypeSuiteObservation,
		},
	}
	assert.Equal(ErrAbsenceConstatNonConforme, inspection.CheckCoherenceSuiteConstats())
}

func TestCheckIncoherenceSuiteConformeEtPointDeControleNonConforme(t *testing.T) {
	assert := require.New(t)
	inspection := Inspection{
		PointsDeControle: []PointDeControle{
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatConforme,
				},
			},
			PointDeControle{
				Publie: true,
				Constat: &Constat{
					Type: TypeConstatObservation,
				},
			},
		},
		Suite: &Suite{
			Type: TypeSuiteAucune,
		},
	}
	assert.Equal(ErrPresenceConstatNonConforme, inspection.CheckCoherenceSuiteConstats())
}
