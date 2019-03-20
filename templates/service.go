package templates

import (
	html "html/template"
	text "text/template"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/fatih/structs"
)

type Config struct {
	Dir     string `default:"templates/"`
	BaseURL string `default:"https://filharmonic.beta.gouv.fr"`
}

type TemplateService struct {
	config Config

	Templates map[TemplateType]*Template
}

func New(config Config) (*TemplateService, error) {
	service := &TemplateService{
		config:    config,
		Templates: make(map[TemplateType]*Template),
	}
	dirEmails := config.Dir + "emails/"
	var err error
	for _, templateName := range emailTemplates {
		template := &Template{}
		template.HTML, err = html.ParseFiles(dirEmails+"layout.html", dirEmails+string(templateName)+".html")
		if err != nil {
			return nil, err
		}
		template.Text, err = text.ParseFiles(dirEmails+"layout.txt", dirEmails+string(templateName)+".txt")
		if err != nil {
			return nil, err
		}
		service.Templates[templateName] = template
	}

	dirOdt := config.Dir + "odt/"
	for _, templateName := range odtTemplates {
		template := &Template{}
		template.Text, err = text.New(string(templateName) + ".fodt").
			Funcs(templateHelpers).
			ParseFiles(dirOdt + string(templateName) + ".fodt")
		if err != nil {
			return nil, err
		}
		service.Templates[templateName] = template
	}
	return service, nil
}

func (s *TemplateService) RenderEmailExpirationDelais(data interface{}) (*models.RenderedTemplate, error) {
	return s.render(TemplateEmailExpirationDelais, data)
}

func (s *TemplateService) RenderEmailNouveauxMessages(data interface{}) (*models.RenderedTemplate, error) {
	return s.render(TemplateEmailNouveauxMessages, data)
}

func (s *TemplateService) RenderEmailRappelEcheances(data interface{}) (*models.RenderedTemplate, error) {
	return s.render(TemplateEmailRappelEcheances, data)
}

func (s *TemplateService) RenderEmailRecapValidation(data interface{}) (*models.RenderedTemplate, error) {
	return s.render(TemplateEmailRecapValidation, data)
}

func (s *TemplateService) RenderODTLettreAnnonce(data interface{}) (*models.RenderedTemplate, error) {
	return s.render(TemplateODTLettreAnnonce, data)
}

func (s *TemplateService) RenderODTLettreSuite(data interface{}) (*models.RenderedTemplate, error) {
	return s.render(TemplateODTLettreSuite, data)
}

func (s *TemplateService) RenderODTRapport(data interface{}) (*models.RenderedTemplate, error) {
	return s.render(TemplateODTRapport, data)
}

func (s *TemplateService) render(template TemplateType, data interface{}) (*models.RenderedTemplate, error) {
	return s.Templates[template].Render(s.addCommonVariables(data))
}

func (s *TemplateService) addCommonVariables(data interface{}) map[string]interface{} {
	output := structs.Map(data)
	output["BaseURL"] = s.config.BaseURL
	return output
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
