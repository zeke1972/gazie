# GAzie TUI - Gestione Aziendale

Un'applicazione TUI (Terminal User Interface) per la gestione aziendale, riscritta completamente in Go utilizzando Bubble Tea e SQLite come database.

## Panoramica

Questa applicazione rappresenta un refactoring completo del sistema ERP GAzie (originariamente in PHP) trasformato in un'applicazione TUI moderna e efficiente utilizzando:

- **Go 1.21+** - Linguaggio di programmazione principale
- **Bubble Tea** - Framework TUI per Go
- **Lipgloss** - Sistema di styling per terminale
- **SQLite3** - Database locale per la persistenza dei dati

## Funzionalità

### Moduli Implementati

1. **📋 Anagrafica Clienti**
   - Visualizzazione lista clienti
   - Aggiunta di nuovi clienti
   - Gestione dati base (nome, codice, città, telefono, email)
   - Navigazione con tastiera

2. **📦 Anagrafica Prodotti**
   - Visualizzazione lista prodotti
   - Aggiunta di nuovi prodotti
   - Gestione prezzi e giacenze
   - Descrizioni prodotti

3. **Database SQLite**
   - Persistenza dati locale
   - Tabelle per clienti e prodotti
   - Dati di esempio precaricati

### Funzionalità Future

- **📄 Ordini Cliente** - Gestione ordini e preventivi
- **🧾 Fatture** - Fatturazione elettronica
- **💰 Contabilità** - Prima nota e registri
- **🏪 Magazzino** - Movimenti di carico/scarico

## Installazione e Compilazione

### Prerequisiti

- Go 1.21 o superiore
- Compilatore C per SQLite3 (di solito preinstallato su Linux/macOS)

### Compilazione

```bash
# Clona il repository
git clone <repository-url>
cd gazie-tui

# Scarica le dipendenze
go mod tidy

# Compila l'applicazione
go build -o gazie-tui

# Esegui l'applicazione
./gazie-tui
```

### Utilizzo

```bash
# Compila ed esegui direttamente
go run main.go
```

## Utilizzo dell'Applicazione

### Controlli Base

- **↑↓** - Naviga tra le voci del menu
- **INVIO** - Seleziona/Conferma
- **ESC** - Torna indietro/Esci
- **N** - Nuovo elemento (nelle viste liste)
- **R** - Ricarica dati
- **Ctrl+C** - Esci dall'applicazione

### Workflow Tipico

1. **Avvio**: L'applicazione si avvia con il menu principale
2. **Navigazione**: Usa ↑↓ per navigare, INVIO per selezionare
3. **Gestione Dati**: 
   - Vai su "Anagrafica Clienti" o "Anagrafica Prodotti"
   - Premi **N** per aggiungere un nuovo elemento
   - Inserisci i dati separati da `|` (pipe)
   - Premi **INVIO** per salvare
4. **Uscita**: Premi **ESC** dal menu principale o **Ctrl+C** in qualsiasi momento

### Esempi di Inserimento Dati

**Cliente**: 
```
Mario Rossi|C001|Roma|06-123456|mario@email.it
```

**Prodotto**:
```
Prodotto示例|P001|15.50|Descrizione del prodotto示例
```

## Architettura

### Struttura del Progetto

```
gazie-tui/
├── main.go                 # Entry point principale
├── go.mod                  # Modulo Go e dipendenze
├── .gitignore             # File ignorati da Git
└── README.md              # Documentazione
```

### Componenti Principali

1. **Models** - Strutture dati (Customer, Product)
2. **Database Layer** - Gestione SQLite e query
3. **UI Layer** - Interfaccia utente con Bubble Tea
4. **Main Model** - Coordinamento tra componenti

### Database Schema

#### Tabella `customers`
- `id` - Chiave primaria
- `code` - Codice cliente univoco
- `name` - Nome/Ragione sociale
- `city` - Città
- `phone` - Numero di telefono
- `email` - Indirizzo email
- `active` - Stato attivo/disattivo
- `created` - Data di creazione

#### Tabella `products`
- `id` - Chiave primaria
- `code` - Codice prodotto univoco
- `name` - Nome prodotto
- `description` - Descrizione
- `price` - Prezzo di vendita
- `stock` - Giacenza attuale
- `min_stock` - Scorte minime
- `active` - Stato attivo/disattivo
- `created` - Data di creazione

## Caratteristiche Tecniche

### Vantaggi del Refactoring

1. **Performance** - Go è molto più veloce di PHP per operazioni CPU-bound
2. **Memoria** - Gestione più efficiente della memoria
3. **Distribuzione** - Compilazione statica, nessuna dipendenza runtime
4. **Portabilità** - Funziona su Windows, Linux, macOS
5. **Manutenibilità** - Codice type-safe con Go

### Caratteristiche TUI

- **Interfaccia moderna** con colori e styling
- **Navigazione intuitiva** con tastiera
- **Responsive design** che si adatta alle dimensioni del terminale
- **Feedback visivo** con status bar e messaggi
- **Gestione errori** con messaggi user-friendly

## Sviluppo

### Estensioni Future

1. **Moduli aggiuntivi** per ordini, fatture, contabilità
2. **Ricerca avanzata** e filtri
3. **Import/Export** dati CSV/Excel
4. **Backup e restore** database
5. **Configurazione** multi-azienda
6. **Reporting** e statistiche

### Personalizzazione

L'applicazione può essere facilmente estesa:

```go
// Aggiungere nuovi stati
type AppState int
const (
    StateMainMenu AppState = iota
    StateCustomers
    StateProducts
    StateOrders        // Nuovo stato
    StateInvoices      // Nuovo stato
)
```

### Testing

```bash
# Esegui i test (quando implementati)
go test ./...

# Verifica codice
go vet ./...

# Formattazione
go fmt ./...
```

## Licenza

Questo progetto è rilasciato sotto licenza MIT. Vedi il file LICENSE per i dettagli.

## Contribuire

1. Fork del repository
2. Crea un branch per la tua feature (`git checkout -b feature/AmazingFeature`)
3. Commit delle modifiche (`git commit -m 'Add some AmazingFeature'`)
4. Push al branch (`git push origin feature/AmazingFeature`)
5. Apri una Pull Request

## Supporto

Per problemi o domande:
- Apri un issue su GitHub
- Consulta la documentazione
- Controlla i log dell'applicazione

---

**Nota**: Questa è una versione di refactoring dell'originale GAzie ERP. Le funzionalità sono in sviluppo attivo e soggette a cambiamenti.