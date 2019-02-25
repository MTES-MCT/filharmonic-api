package templates

import (
	"bytes"
	html "html/template"
	text "text/template"

	"github.com/MTES-MCT/filharmonic-api/models"
)

type Config struct {
	Dir string `default:"templates/templates/"`
}

type TemplateService struct {
	config Config

	emailNouveauxMessagesTemplate *html.Template
	lettreAnnonceTemplate         *text.Template
	lettreSuiteTemplate           *text.Template
	rapportTemplate               *text.Template
}

var (
	templateHelpers = text.FuncMap{
		"type_suite": func(value models.TypeSuite) string {
			switch value {
			case models.TypeSuiteAucune:
				return "Aucune"
			case models.TypeSuiteObservation:
				return "Observation ou non conformités à traiter par courrier"
			case models.TypeSuitePropositionMiseEnDemeure:
				return "Proposition de suites administratives"
			case models.TypeSuitePropositionRenforcement:
				return "Proposition de renforcement, modification ou mise à jour des prescription"
			case models.TypeSuiteAutre:
				return "Autre"
			default:
				return string(value)
			}
		},
		"type_constat": func(value models.TypeConstat) string {
			switch value {
			case models.TypeConstatConforme:
				return "Conforme"
			case models.TypeConstatNonConforme:
				return "Non conforme"
			case models.TypeConstatObservation:
				return "Observation"
			default:
				return string(value)
			}
		},
		"type_inspection": func(value models.TypeInspection) string {
			switch value {
			case models.TypeApprofondi:
				return "Approfondi"
			case models.TypeCourant:
				return "Courant"
			case models.TypePonctuel:
				return "Ponctuel"
			default:
				return string(value)
			}
		},
		"ouinon": func(value bool) string {
			if value {
				return "Oui"
			}
			return "Non"
		},
		"origine_inspection": func(value models.OrigineInspection) string {
			switch value {
			case models.OrigineCirconstancielle:
				return "Circonstancielle"
			case models.OriginePlanControle:
				return "Plan de contrôle"
			default:
				return string(value)
			}
		},
		"circonstance_inspection": func(value models.CirconstanceInspection) string {
			switch value {
			case models.CirconstanceAutre:
				return "Autre"
			case models.CirconstanceIncident:
				return "Incident"
			case models.CirconstancePlainte:
				return "Plainte"
			default:
				return string(value)
			}
		},
		"regime_etablissement": func(value models.RegimeEtablissement) string {
			switch value {
			case models.RegimeAucun:
				return "Aucun"
			case models.RegimeAutorisation:
				return "Autorisation"
			case models.RegimeDeclaration:
				return "Déclaration"
			case models.RegimeEnregistrement:
				return "Enregistrement"
			default:
				return string(value)
			}
		},
		"add": func(a int, b int) int {
			return a + b
		},
	}
)

func New(config Config) (*TemplateService, error) {
	service := &TemplateService{
		config: config,
	}
	var err error
	service.emailNouveauxMessagesTemplate, err = html.ParseFiles(config.Dir + "email-new-messages.html")
	if err != nil {
		return nil, err
	}

	service.lettreAnnonceTemplate, err = text.New("modele_de_lettre_annonce_inspection.fodt").
		Funcs(templateHelpers).
		ParseFiles(config.Dir + "modele_de_lettre_annonce_inspection.fodt")
	if err != nil {
		return nil, err
	}

	service.lettreSuiteTemplate, err = text.New("modele_de_lettre_suite_inspection.fodt").
		Funcs(templateHelpers).
		ParseFiles(config.Dir + "modele_de_lettre_suite_inspection.fodt")
	if err != nil {
		return nil, err
	}

	service.rapportTemplate, err = text.New("modele_de_rapport_inspection.fodt").
		Funcs(templateHelpers).
		ParseFiles(config.Dir + "modele_de_rapport_inspection.fodt")
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *TemplateService) RenderHTMLEmailNouveauxMessages(data interface{}) (string, error) {
	return s.renderHTMLTemplate(s.emailNouveauxMessagesTemplate, data)
}

func (s *TemplateService) RenderLettreAnnonce(data interface{}) (string, error) {
	return s.renderTextTemplate(s.lettreAnnonceTemplate, data)
}

func (s *TemplateService) RenderLettreSuite(data interface{}) (string, error) {
	return s.renderTextTemplate(s.lettreSuiteTemplate, data)
}

func (s *TemplateService) RenderRapport(data interface{}) (string, error) {
	return s.renderTextTemplate(s.rapportTemplate, data)
}

func (s *TemplateService) renderHTMLTemplate(tmpl *html.Template, data interface{}) (string, error) {
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func (s *TemplateService) renderTextTemplate(tmpl *text.Template, data interface{}) (string, error) {
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
