package seeds

import (
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/types"
)

/*
Initialise la base de données avec un jeu de test minimal, utilisable dans les tests.

Attention à ne pas préciser les Id sans quoi les séquences de clés primaires ne sont pas incrémentées.
*/
// nolint: gocyclo
func SeedsTestDB(db *pg.DB) error {
	users := []interface{}{
		&models.User{
			// Id:      1,
			Prenom:  "Michel",
			Nom:     "Exploitant1",
			Email:   "exploitant1@filharmonic.com",
			Profile: models.ProfilExploitant,
		},
		&models.User{
			// Id:      2,
			Prenom:  "Bernard",
			Nom:     "Exploitant2",
			Email:   "exploitant2@filharmonic.com",
			Profile: models.ProfilExploitant,
		},
		&models.User{
			// Id:      3,
			Prenom:  "Alain",
			Nom:     "Champion",
			Email:   "inspecteur1@filharmonic.com",
			Profile: models.ProfilInspecteur,
		},
		&models.User{
			// Id:      4,
			Prenom:  "Corine",
			Nom:     "Dupont",
			Email:   "inspecteur2@filharmonic.com",
			Profile: models.ProfilInspecteur,
		},
		&models.User{
			// Id:      5,
			Prenom:  "Bernard",
			Nom:     "Mars",
			Email:   "inspecteur3@filharmonic.com",
			Profile: models.ProfilInspecteur,
		},
		&models.User{
			// Id:      6,
			Prenom:  "Albert",
			Nom:     "Approbe",
			Email:   "approbateur1@filharmonic.com",
			Profile: models.ProfilApprobateur,
		},
		&models.User{
			// Id:      7,
			Prenom:  "Gilbert",
			Nom:     "Approbe",
			Email:   "approbateur2@filharmonic.com",
			Profile: models.ProfilApprobateur,
		},
		&models.User{
			// Id:      8,
			Prenom:  "Hubert",
			Nom:     "Inspecte",
			Email:   "inspecteur4@filharmonic.com",
			Profile: models.ProfilInspecteur,
		},
	}
	err := db.Insert(users...)
	if err != nil {
		return err
	}

	themes := []interface{}{
		&models.Theme{
			// Id:      1,
			Nom: "Rejets dans l'eau",
		},
		&models.Theme{
			// Id:      2,
			Nom: "Rejets dans l'air",
		},
		&models.Theme{
			// Id:      3,
			Nom: "Sûreté",
		},
		&models.Theme{
			// Id:      4,
			Nom: "Produits chimiques",
		},
		&models.Theme{
			// Id:      5,
			Nom: "Incendie",
		},
		&models.Theme{
			// Id:      6,
			Nom: "COV",
		},
	}
	err = db.Insert(themes...)
	if err != nil {
		return err
	}

	departements := []interface{}{
		&models.Departement{
			// Id: 1
			CodeInsee:       "75",
			Nom:             "Paris",
			Charniere:       "de ",
			Region:          "Île-de-France",
			CharniereRegion: "d'",
		},
		&models.Departement{
			// Id: 2
			CodeInsee:       "44",
			Nom:             "Loire-Atlantique",
			Charniere:       "de la ",
			Region:          "Pays de la Loire",
			CharniereRegion: "des ",
		},
		&models.Departement{
			// Id: 3
			CodeInsee:       "69",
			Nom:             "Rhône",
			Charniere:       "du ",
			Region:          "Auvergne-Rhône-Alpes",
			CharniereRegion: "d'",
		},
		&models.Departement{
			// Id: 4
			CodeInsee:       "63",
			Nom:             "Puy-de-Dôme",
			Charniere:       "du ",
			Region:          "Auvergne-Rhône-Alpes",
			CharniereRegion: "d'",
		},
		&models.Departement{
			// Id: 5
			CodeInsee:       "01",
			Nom:             "Ain",
			Charniere:       "de l'",
			Region:          "Auvergne-Rhône-Alpes",
			CharniereRegion: "d'",
		},
	}
	err = db.Insert(departements...)
	if err != nil {
		return err
	}

	etablissements := []interface{}{
		&models.Etablissement{
			// Id:      1,
			S3IC:          "1234",
			Nom:           "Nom 1",
			Raison:        "Raison sociale",
			Adresse1:      "1 rue des fleurs",
			Adresse2:      "",
			CodePostal:    "75000",
			Commune:       "Paris",
			Regime:        models.RegimeAutorisation,
			DepartementId: 1,
		},
		&models.Etablissement{
			// Id:      2,
			S3IC:          "451267",
			Nom:           "Nom 2",
			Raison:        "Raison sociale 2",
			Adresse1:      "1 rue des plantes",
			Adresse2:      "parcelle 207",
			CodePostal:    "44000",
			Commune:       "Nantes",
			Regime:        models.RegimeDeclaration,
			DepartementId: 2,
		},
		&models.Etablissement{
			// Id:      3,
			S3IC:          "3335655",
			Nom:           "Nom 3",
			Raison:        "Raison sociale 3",
			Adresse1:      "1 rue des cordeliers",
			Adresse2:      "",
			CodePostal:    "69000",
			Commune:       "Lyon",
			Regime:        models.RegimeAucun,
			DepartementId: 3,
		},
		&models.Etablissement{
			// Id:      4,
			S3IC:          "4444213",
			Nom:           "Nom 4",
			Raison:        "Raison sociale 4",
			Adresse1:      "1 place de l'église",
			Adresse2:      "",
			CodePostal:    "63000",
			Commune:       "Clermont-Ferrand",
			Regime:        models.RegimeEnregistrement,
			DepartementId: 4,
		},
	}
	err = db.Insert(etablissements...)
	if err != nil {
		return err
	}

	etablissementToExploitants := []interface{}{
		&models.EtablissementToExploitant{
			EtablissementId: 1,
			UserId:          1,
		},
		&models.EtablissementToExploitant{
			EtablissementId: 2,
			UserId:          1,
		},
		&models.EtablissementToExploitant{
			EtablissementId: 3,
			UserId:          2,
		},
		&models.EtablissementToExploitant{
			EtablissementId: 4,
			UserId:          2,
		},
	}
	err = db.Insert(etablissementToExploitants...)
	if err != nil {
		return err
	}

	suites := []interface{}{
		&models.Suite{
			// Id: 1,
			Type:     models.TypeSuiteObservation,
			Synthese: "Observations à traiter",
		},
		&models.Suite{
			// Id: 2,
			Type: models.TypeSuitePropositionMiseEnDemeure,
		},
		&models.Suite{
			// Id: 3,
			Type: models.TypeSuitePropositionMiseEnDemeure,
		},
	}
	err = db.Insert(suites...)
	if err != nil {
		return err
	}

	rapports := []interface{}{
		&models.Rapport{
			// Id:   1,
			Nom:       "rapport.pdf",
			AuteurId:  6,
			StorageId: "rapport-1234",
		},
	}
	err = db.Insert(rapports...)
	if err != nil {
		return err
	}

	canevas := []interface{}{
		&models.Canevas{
			// Id:   1,
			Nom:      "test",
			AuteurId: 3,
			Data: models.CanevasData{
				PointsDeControle: []models.CanevasPointDeControle{
					models.CanevasPointDeControle{
						Sujet: "COV",
						ReferencesReglementaires: []string{
							"Article 5 du 01/02/1995",
						},
					},
					models.CanevasPointDeControle{
						Sujet: "Rejets dans l'air",
						ReferencesReglementaires: []string{
							"Article 1 du 01/02/2003",
							"Article 2 du 01/02/2002",
						},
						Message: "Donnez-moi le document sur l'air.",
					},
					models.CanevasPointDeControle{
						Sujet: "Rejets dans l'eau",
						ReferencesReglementaires: []string{
							"Article 13 du 01/02/2006",
						},
						Message: "Donnez-moi le document sur l'eau.",
					},
				},
			},
			DataVersion: 1,
			CreatedAt:   util.DateTime("2019-03-11T08:30:45"),
		},
	}
	err = db.Insert(canevas...)
	if err != nil {
		return err
	}

	inspections := []interface{}{
		&models.Inspection{ // (non cohérent) en cours, avec suite, 1 point de contrôle publié avec constat, 1 point de contrôle non publié
			// Id:   1,
			Date: util.Date("2018-09-01"),
			Type: models.TypeApprofondi,
			Themes: []string{
				"Produits chimiques",
			},
			Annonce:         true,
			Origine:         models.OriginePlanControle,
			Etat:            models.EtatEnCours,
			Contexte:        "Emissions de NOx dépassant les seuils le 2/04/2005",
			EtablissementId: 1,
			SuiteId:         1,
		},
		&models.Inspection{ // (non cohérent) en préparation, pas de suite, 2 points de contrôle publiés dont un avec constat
			// Id:   2,
			Date: util.Date("2018-11-15"),
			Type: models.TypeCourant,
			Themes: []string{
				"Rejets dans l'eau",
				"Rejets dans l'air",
			},
			Annonce:         true,
			Origine:         models.OriginePlanControle,
			Etat:            models.EtatPreparation,
			Contexte:        "Incident cuve gaz le 30/12/2017",
			EtablissementId: 3,
		},
		&models.Inspection{ // en attente validation, avec suite, 1 point de contrôle publié avec constat
			// Id:   3,
			Date: util.Date("2018-11-15"),
			Type: models.TypeApprofondi,
			Themes: []string{
				"Sanitaire",
			},
			Annonce:         true,
			Origine:         models.OriginePlanControle,
			Etat:            models.EtatAttenteValidation,
			Contexte:        "Viande avariée",
			EtablissementId: 4,
			SuiteId:         2,
		},
		&models.Inspection{ // en cours, pas de suite, 2 points de contrôle publiés dont un sans constat
			// Id:   4,
			Date: util.Date("2019-01-15"),
			Type: models.TypeApprofondi,
			Themes: []string{
				"Sanitaire",
			},
			Annonce:         true,
			Origine:         models.OriginePlanControle,
			Etat:            models.EtatEnCours,
			Contexte:        "Inspection en cours",
			EtablissementId: 4,
		},
		&models.Inspection{ // en attente validation, avec suite, 1 point de contrôle publié avec constat
			// Id:   5,
			Date: util.Date("2019-02-26"),
			Type: models.TypeApprofondi,
			Themes: []string{
				"Rejets dans l'air",
			},
			Annonce: true,
			Origine: models.OriginePlanControle,
			Etat:    models.EtatTraitementNonConformites,
			DateValidation: types.NullTime{
				Time: util.DateTime("2019-02-28T12:00:00"),
			},
			Contexte:        "Rejets CO2 de l'usine",
			EtablissementId: 4,
			SuiteId:         3,
			RapportId:       1,
		},
	}
	err = db.Insert(inspections...)
	if err != nil {
		return err
	}

	inspectionToInspecteurs := []interface{}{
		&models.InspectionToInspecteur{
			InspectionId: 1,
			UserId:       3,
		},
		&models.InspectionToInspecteur{
			InspectionId: 2,
			UserId:       3,
		},
		&models.InspectionToInspecteur{
			InspectionId: 3,
			UserId:       3,
		},
		&models.InspectionToInspecteur{
			InspectionId: 3,
			UserId:       4,
		},
		&models.InspectionToInspecteur{
			InspectionId: 4,
			UserId:       3,
		},
		&models.InspectionToInspecteur{
			InspectionId: 1,
			UserId:       4,
		},
		&models.InspectionToInspecteur{
			InspectionId: 2,
			UserId:       5,
		},
		&models.InspectionToInspecteur{
			InspectionId: 5,
			UserId:       3,
		},
	}
	err = db.Insert(inspectionToInspecteurs...)
	if err != nil {
		return err
	}

	userToFavoris := []interface{}{
		&models.UserToFavori{
			InspectionId: 1,
			UserId:       3,
		},
		&models.UserToFavori{
			InspectionId: 1,
			UserId:       1,
		},
	}
	err = db.Insert(userToFavoris...)
	if err != nil {
		return err
	}

	constats := []interface{}{
		&models.Constat{
			// Id: 1,
			Type:      models.TypeConstatObservation,
			Remarques: "Il manque des choses",
		},
		&models.Constat{
			// Id: 2,
			Type:      models.TypeConstatConforme,
			Remarques: "RAS",
		},
		&models.Constat{
			// Id: 3,
			Type:      models.TypeConstatNonConforme,
			Remarques: "Ne respecte pas la réglementation",
		},
		&models.Constat{
			// Id: 4,
			Type:               models.TypeConstatNonConforme,
			Remarques:          "Ne respecte pas la réglementation",
			DelaiNombre:        7,
			DelaiUnite:         "jours",
			EcheanceResolution: util.Date("2019-03-06"),
		},
	}
	err = db.Insert(constats...)
	if err != nil {
		return err
	}

	pointsDeControle := []interface{}{
		&models.PointDeControle{
			// Id:    1,
			Sujet: "Mesure des émissions atmosphériques canalisées par un organisme extérieur",
			ReferencesReglementaires: []string{
				"Article 3.2.3. de l'arrêté préfectoral du 28 juin 2017",
				"Article 3.2.8. de l'arrêté préfectoral du 28 juin 2017",
				"Article 8.2.1.2. de l'arrêté préfectoral du 28 juin 2017",
			},
			Publie:       true,
			InspectionId: 1,
			ConstatId:    1,
		},
		&models.PointDeControle{
			// Id:    2,
			Sujet: "Autosurveillance des émissions canalisées de COV",
			ReferencesReglementaires: []string{
				"Article 8.2.1.1. de l'arrêté préfectoral du 28 juin 2017",
			},
			Publie:       false,
			InspectionId: 1,
		},
		&models.PointDeControle{
			// Id:    3,
			Sujet: "Eau - Air",
			ReferencesReglementaires: []string{
				"Article 1 de l'Arrêté ministériel du 28 avril 2014",
			},
			Publie:       true,
			InspectionId: 2,
		},
		&models.PointDeControle{
			// Id:    4,
			Sujet: "Santé",
			ReferencesReglementaires: []string{
				"Article 1 de l'Arrêté ministériel du 28 avril 2014",
			},
			Publie:       true,
			InspectionId: 3,
			ConstatId:    2,
		},
		&models.PointDeControle{
			// Id:    5,
			Sujet: "Santé 1",
			ReferencesReglementaires: []string{
				"Article 1 de l'Arrêté ministériel du 28 avril 2014",
			},
			Publie:       true,
			InspectionId: 4,
			ConstatId:    3,
			Order:        1,
		},
		&models.PointDeControle{
			// Id:    6,
			Sujet: "Santé 2",
			ReferencesReglementaires: []string{
				"Article 1 de l'Arrêté ministériel du 228 avril 2014",
			},
			Publie:       true,
			InspectionId: 4,
			Order:        2,
		},
		&models.PointDeControle{
			// Id:    7,
			Sujet: "Rejets Air",
			ReferencesReglementaires: []string{
				"Article 1 de l'Arrêté ministériel du 28 avril 2015",
			},
			Publie:       true,
			InspectionId: 5,
			ConstatId:    4,
		},
	}
	err = db.Insert(pointsDeControle...)
	if err != nil {
		return err
	}

	commentaires := []interface{}{
		&models.Commentaire{
			// Id:           1,
			Message:      "Attention à l'article 243.",
			Date:         util.DateTime("2018-11-14T08:50:00"),
			AuteurId:     3,
			InspectionId: 1,
		},
		&models.Commentaire{
			// Id:           2,
			Message:      "L'article 843 s'applique également.",
			Date:         util.DateTime("2018-11-16T16:50:00"),
			AuteurId:     4,
			InspectionId: 1,
		},
		&models.Commentaire{
			// Id:           3,
			Message:      "Attention au précédent contrôle.",
			Date:         util.DateTime("2018-11-18T16:50:00"),
			AuteurId:     3,
			InspectionId: 2,
		},
	}
	err = db.Insert(commentaires...)
	if err != nil {
		return err
	}

	messages := []interface{}{
		&models.Message{
			// Id:                1,
			Message:           "Auriez-vous l'obligeance de me fournir le document approprié ?",
			Date:              util.DateTime("2018-11-14T08:50:00"),
			Lu:                true,
			Interne:           false,
			AuteurId:          3,
			PointDeControleId: 1,
		},
		&models.Message{
			// Id:                2,
			Message:           "Voici le document.",
			Date:              util.DateTime("2018-11-16T16:50:00"),
			Lu:                true,
			Interne:           false,
			AuteurId:          1,
			PointDeControleId: 1,
		},
		&models.Message{
			// Id:                3,
			Message:           "Attention au précédent contrôle.",
			Date:              util.DateTime("2018-11-20T16:50:00"),
			Lu:                false,
			Interne:           true,
			AuteurId:          3,
			PointDeControleId: 1,
		},
		&models.Message{
			// Id:                4,
			Message:           "Merci de me fournir le document.",
			Date:              util.DateTime("2018-11-21T16:50:00"),
			Lu:                false,
			Interne:           false,
			AuteurId:          4,
			PointDeControleId: 2,
		},
		&models.Message{
			// Id:                5,
			Message:           "Auriez-vous l'obligeance de me fournir une photo de la cuve ?",
			Date:              util.DateTime("2018-11-18T17:50:00"),
			Lu:                true,
			Interne:           false,
			AuteurId:          3,
			PointDeControleId: 3,
		},
		&models.Message{
			// Id:                6,
			Message:           "Il manque un document.",
			Date:              util.DateTime("2018-11-26T16:50:00"),
			Lu:                false,
			Interne:           false,
			AuteurId:          3,
			PointDeControleId: 1,
		},
		&models.Message{
			// Id:                7,
			Message:           "Voici la photo de la cuve.",
			Date:              util.DateTime("2018-11-27T16:50:00"),
			Lu:                false,
			Interne:           false,
			AuteurId:          2,
			PointDeControleId: 3,
		},
	}
	err = db.Insert(messages...)
	if err != nil {
		return err
	}

	piecesJointes := []interface{}{
		&models.PieceJointe{
			// Id:                1,
			Nom:       "photo-cuve.pdf",
			Type:      "application/pdf",
			Taille:    7945,
			MessageId: 2,
			StorageId: "photo-cuve.pdf",
			AuteurId:  1,
		},
		&models.PieceJointe{
			// Id:                2,
			Nom:       "photo-cuve-2.pdf",
			Type:      "application/pdf",
			Taille:    2262000,
			StorageId: "1234567890",
			AuteurId:  1,
		},
		&models.PieceJointe{
			// Id:                3,
			Nom:           "article-loi.pdf",
			Type:          "application/pdf",
			Taille:        2262000,
			StorageId:     "3456789054",
			CommentaireId: 2,
			AuteurId:      3,
		},
		&models.PieceJointe{
			// Id:                4,
			Nom:       "article-loi-v2.pdf",
			Type:      "application/pdf",
			Taille:    10000000,
			StorageId: "4567890543",
			AuteurId:  3,
		},
		&models.PieceJointe{
			// Id:                5,
			Nom:       "image-cuve.jpg",
			Type:      "image/jpeg",
			Taille:    53553,
			MessageId: 2,
			StorageId: "cuve.jpg",
			AuteurId:  1,
		},
		&models.PieceJointe{
			// Id:                5,
			Nom:       "rapport-cuve.odt",
			Type:      "application/vnd.oasis.opendocument.text",
			Taille:    6737,
			MessageId: 2,
			StorageId: "rapport-cuve.odt",
			AuteurId:  1,
		},
	}
	err = db.Insert(piecesJointes...)
	if err != nil {
		return err
	}

	evenements := []interface{}{
		&models.Evenement{
			// Id:        1,
			Type:         models.EvenementCreationMessage,
			AuteurId:     3,
			InspectionId: 1,
			CreatedAt:    util.DateTime("2018-09-16T14:00:00"),
			Data: map[string]interface{}{
				"message_id":           1,
				"point_de_controle_id": 1,
			},
		},
		&models.Evenement{
			// Id:        2,
			Type:         models.EvenementCommentaireGeneral,
			AuteurId:     3,
			InspectionId: 2,
			CreatedAt:    util.DateTime("2018-11-14T08:50:00"),
			Data: map[string]interface{}{
				"message_d": 2,
			},
		},
		&models.Evenement{
			// Id:        3,
			Type:         models.EvenementModificationInspection,
			AuteurId:     3,
			InspectionId: 2,
			CreatedAt:    util.DateTime("2018-12-07T18:42:20"),
		},
		&models.Evenement{
			// Id:        4,
			Type:         models.EvenementLectureMessage,
			AuteurId:     3,
			InspectionId: 3,
			CreatedAt:    util.DateTime("2019-01-07T18:42:20"),
		},
	}

	err = db.Insert(evenements...)
	if err != nil {
		return err
	}

	notifications := []interface{}{
		&models.Notification{
			// Id:       1,
			Lue:            false,
			EvenementId:    1,
			DestinataireId: 3,
		},
		&models.Notification{
			// Id:       2,
			Lue:            false,
			EvenementId:    2,
			DestinataireId: 3,
		},
		&models.Notification{
			// Id:       3,
			Lue:            false,
			EvenementId:    3,
			DestinataireId: 3,
		},
	}

	return db.Insert(notifications...)
}
