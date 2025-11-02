package dialog

import (
	fdialog "github.com/sqweek/dialog"
)

// DialogType rappresenta il tipo di dialog
type DialogType string

const (
	// DialogTypeInfo mostra un dialog informativo
	DialogTypeInfo DialogType = "info"
	// DialogTypeError mostra un dialog di errore
	DialogTypeError DialogType = "error"
	// DialogTypeQuestion mostra un dialog con domanda (Yes/No)
	DialogTypeQuestion DialogType = "question"
	// DialogTypeFileOpen mostra un dialog per aprire file
	DialogTypeFileOpen DialogType = "file-open"
	// DialogTypeFileSave mostra un dialog per salvare file
	DialogTypeFileSave DialogType = "file-save"
	// DialogTypeDirSelect mostra un dialog per selezionare directory
	DialogTypeDirSelect DialogType = "dir-select"
)

// DialogConfig contiene la configurazione per un dialog
type DialogConfig struct {
	Type        DialogType
	Title       string
	Message     string
	DefaultPath string
	Filters     []FileFilter
}

// FileFilter rappresenta un filtro per file
type FileFilter struct {
	Description string
	Pattern     string
}

// DialogResult contiene il risultato di un dialog
type DialogResult struct {
	OK       bool   // true se l'utente ha cliccato OK/Yes
	FilePath string // percorso del file selezionato (per file dialogs)
	Error    error  // eventuale errore
}

// Show mostra un dialog e restituisce il risultato
func Show(config DialogConfig) DialogResult {
	switch config.Type {
	case DialogTypeInfo:
		return showInfo(config)
	case DialogTypeError:
		return showError(config)
	case DialogTypeQuestion:
		return showQuestion(config)
	case DialogTypeFileOpen:
		return showFileOpen(config)
	case DialogTypeFileSave:
		return showFileSave(config)
	case DialogTypeDirSelect:
		return showDirSelect(config)
	default:
		return DialogResult{OK: false, Error: nil}
	}
}

// showInfo mostra un dialog informativo
func showInfo(config DialogConfig) DialogResult {
	msg := config.Message
	if msg == "" {
		msg = "Information"
	}
	title := config.Title
	if title == "" {
		title = "Info"
	}
	fdialog.Message("%s", msg).Title(title).Info()
	return DialogResult{OK: true}
}

// showError mostra un dialog di errore
func showError(config DialogConfig) DialogResult {
	msg := config.Message
	if msg == "" {
		msg = "Error"
	}
	title := config.Title
	if title == "" {
		title = "Error"
	}
	fdialog.Message("%s", msg).Title(title).Error()
	return DialogResult{OK: true}
}

// showQuestion mostra un dialog con domanda Yes/No
func showQuestion(config DialogConfig) DialogResult {
	msg := config.Message
	if msg == "" {
		msg = "Are you sure?"
	}
	title := config.Title
	if title == "" {
		title = "Question"
	}
	ok := fdialog.Message("%s", msg).Title(title).YesNo()
	return DialogResult{OK: ok}
}

// showFileOpen mostra un dialog per aprire file
func showFileOpen(config DialogConfig) DialogResult {
	dlg := fdialog.File().Title(config.Title)

	// Imposta la directory di default se specificata
	if config.DefaultPath != "" {
		dlg = dlg.SetStartDir(config.DefaultPath)
	}

	// Aggiungi filtri se specificati
	for _, filter := range config.Filters {
		dlg = dlg.Filter(filter.Description, filter.Pattern)
	}

	filename, err := dlg.Load()
	if err != nil {
		// L'utente ha annullato
		return DialogResult{OK: false, Error: err}
	}

	return DialogResult{OK: true, FilePath: filename}
}

// showFileSave mostra un dialog per salvare file
func showFileSave(config DialogConfig) DialogResult {
	dlg := fdialog.File().Title(config.Title)

	// Imposta la directory di default se specificata
	if config.DefaultPath != "" {
		dlg = dlg.SetStartDir(config.DefaultPath)
	}

	// Aggiungi filtri se specificati
	for _, filter := range config.Filters {
		dlg = dlg.Filter(filter.Description, filter.Pattern)
	}

	filename, err := dlg.Save()
	if err != nil {
		// L'utente ha annullato
		return DialogResult{OK: false, Error: err}
	}

	return DialogResult{OK: true, FilePath: filename}
}

// showDirSelect mostra un dialog per selezionare una directory
func showDirSelect(config DialogConfig) DialogResult {
	dlg := fdialog.Directory().Title(config.Title)

	// Imposta la directory di default se specificata
	if config.DefaultPath != "" {
		dlg = dlg.SetStartDir(config.DefaultPath)
	}

	dirname, err := dlg.Browse()
	if err != nil {
		// L'utente ha annullato
		return DialogResult{OK: false, Error: err}
	}

	return DialogResult{OK: true, FilePath: dirname}
}

// Info mostra un dialog informativo (shortcut)
func Info(title, message string) {
	Show(DialogConfig{
		Type:    DialogTypeInfo,
		Title:   title,
		Message: message,
	})
}

// Error mostra un dialog di errore (shortcut)
func Error(title, message string) {
	Show(DialogConfig{
		Type:    DialogTypeError,
		Title:   title,
		Message: message,
	})
}

// Question mostra un dialog con domanda Yes/No (shortcut)
func Question(title, message string) bool {
	result := Show(DialogConfig{
		Type:    DialogTypeQuestion,
		Title:   title,
		Message: message,
	})
	return result.OK
}

// FileOpen mostra un dialog per aprire file (shortcut)
func FileOpen(title string, filters ...FileFilter) (string, error) {
	result := Show(DialogConfig{
		Type:    DialogTypeFileOpen,
		Title:   title,
		Filters: filters,
	})
	return result.FilePath, result.Error
}

// FileSave mostra un dialog per salvare file (shortcut)
func FileSave(title string, filters ...FileFilter) (string, error) {
	result := Show(DialogConfig{
		Type:    DialogTypeFileSave,
		Title:   title,
		Filters: filters,
	})
	return result.FilePath, result.Error
}

// DirSelect mostra un dialog per selezionare directory (shortcut)
func DirSelect(title string) (string, error) {
	result := Show(DialogConfig{
		Type:  DialogTypeDirSelect,
		Title: title,
	})
	return result.FilePath, result.Error
}
