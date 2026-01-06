package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    _ "github.com/mattn/go-sqlite3"

    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Database paths
var (
    AppDataDir = os.Getenv("HOME") + "/.gazie-tui"
    DBPath     = filepath.Join(AppDataDir, "gazie.db")
)

// Models
type Customer struct {
    ID       int
    Code     string
    Name     string
    City     string
    Phone    string
    Email    string
    Active   bool
    Created  time.Time
}

type Product struct {
    ID           int
    Code         string
    Name         string
    Description  string
    Price        float64
    Stock        int
    MinStock     int
    Active       bool
    Created      time.Time
}

// AppState represents the current state of the application
type AppState int

const (
    StateMainMenu AppState = iota
    StateCustomers
    StateProducts
    StateCustomerForm
    StateProductForm
)

// MainModel represents the main application model
type MainModel struct {
    db         *sql.DB
    currentApp AppState
    selected   int
    menuItems  []string
    customers  []Customer
    products   []Product
    formData   string
    formMode   string
    status     string
    width      int
    height     int
    err        error
}

// NewMainModel creates a new main model
func NewMainModel(db *sql.DB) MainModel {
    // Initialize database with sample data
    initDatabase(db)
    
    // Load data
    customers := loadCustomers(db)
    products := loadProducts(db)
    
    return MainModel{
        db:         db,
        currentApp: StateMainMenu,
        selected:   0,
        menuItems:  []string{"📋 Anagrafica Clienti", "📦 Anagrafica Prodotti", "❌ Esci"},
        customers:  customers,
        products:   products,
        formData:   "",
        formMode:   "",
        status:     "Benvenuto in GAzie TUI",
    }
}

// Init initializes the model
func (m MainModel) Init() tea.Cmd {
    return nil
}

// Update handles messages and updates the model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKey(msg)
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}

// handleKey handles keyboard input
func (m MainModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.Type {
    case tea.KeyCtrlC, tea.KeyEsc:
        if m.currentApp == StateCustomerForm || m.currentApp == StateProductForm {
            m.currentApp = StateMainMenu
            m.formData = ""
            m.status = "Annullato"
            return m, nil
        }
        return m, tea.Quit
    case tea.KeyUp:
        if m.selected > 0 {
            m.selected--
        }
    case tea.KeyDown:
        maxItems := len(m.menuItems) - 1
        if m.currentApp == StateCustomers {
            maxItems = len(m.customers) - 1
        } else if m.currentApp == StateProducts {
            maxItems = len(m.products) - 1
        }
        if m.selected < maxItems {
            m.selected++
        }
    case tea.KeyEnter:
        return m.handleEnter()
    case tea.KeyRunes:
        if msg.String() == "r" {
            // Refresh data
            m.customers = loadCustomers(m.db)
            m.products = loadProducts(m.db)
            m.status = "Dati aggiornati"
        } else if msg.String() == "n" {
            if m.currentApp == StateCustomers {
                m.currentApp = StateCustomerForm
                m.formMode = "new"
                m.formData = ""
                m.status = "Nuovo cliente"
            } else if m.currentApp == StateProducts {
                m.currentApp = StateProductForm
                m.formMode = "new"
                m.formData = ""
                m.status = "Nuovo prodotto"
            }
        }
    case tea.KeyBackspace:
        if m.currentApp == StateCustomerForm || m.currentApp == StateProductForm {
            if len(m.formData) > 0 {
                m.formData = m.formData[:len(m.formData)-1]
            }
        }
    case tea.KeySpace:
        if m.currentApp == StateCustomerForm || m.currentApp == StateProductForm {
            m.formData += " "
        }
    default:
        // Handle printable characters
        if len(msg.String()) == 1 && (m.currentApp == StateCustomerForm || m.currentApp == StateProductForm) {
            // Only accept printable ASCII characters and some special chars for forms
            char := msg.String()
            if (char >= "a" && char <= "z") || (char >= "A" && char <= "Z") || 
               (char >= "0" && char <= "9") || char == "|" || char == "." || 
               char == "@" || char == "-" || char == "_" {
                m.formData += char
            }
        }
    }
    return m, nil
}

// handleEnter handles enter key
func (m MainModel) handleSelection() (tea.Model, tea.Cmd) {
    switch m.selected {
    case 0:
        m.currentApp = StateCustomers
        m.status = "Gestione Clienti"
    case 1:
        m.currentApp = StateProducts
        m.status = "Gestione Prodotti"
    case 2:
        return m, tea.Quit
    }
    return m, nil
}

// handleEnter handles enter key based on current state
func (m MainModel) handleEnter() (tea.Model, tea.Cmd) {
    switch m.currentApp {
    case StateMainMenu:
        return m.handleSelection()
    case StateCustomerForm:
        return m.saveCustomer()
    case StateProductForm:
        return m.saveProduct()
    }
    return m, nil
}

// saveCustomer saves a customer
func (m MainModel) saveCustomer() (tea.Model, tea.Cmd) {
    fields := strings.Split(m.formData, "|")
    if len(fields) < 4 {
        m.status = "Errore: insufficienti dati (Nome|Code|Città|Telefono)"
        return m, nil
    }
    
    // Generate new code if empty
    if fields[1] == "" {
        fields[1] = fmt.Sprintf("C%03d", len(m.customers)+1)
    }
    
    query := `INSERT INTO customers (code, name, city, phone, email, active) VALUES (?, ?, ?, ?, ?, 1)`
    _, err := m.db.Exec(query, fields[1], fields[0], fields[2], fields[3], fields[4])
    if err != nil {
        m.status = fmt.Sprintf("Errore: %v", err)
        return m, nil
    }
    
    m.customers = loadCustomers(m.db)
    m.currentApp = StateCustomers
    m.formData = ""
    m.status = "Cliente salvato con successo"
    return m, nil
}

// saveProduct saves a product
func (m MainModel) saveProduct() (tea.Model, tea.Cmd) {
    fields := strings.Split(m.formData, "|")
    if len(fields) < 4 {
        m.status = "Errore: insufficienti dati (Nome|Code|Prezzo|Descrizione)"
        return m, nil
    }
    
    // Parse price
    price, err := strconv.ParseFloat(fields[2], 64)
    if err != nil {
        m.status = "Errore: prezzo non valido"
        return m, nil
    }
    
    // Generate new code if empty
    if fields[1] == "" {
        fields[1] = fmt.Sprintf("P%03d", len(m.products)+1)
    }
    
    query := `INSERT INTO products (code, name, description, price, stock, min_stock, active) VALUES (?, ?, ?, ?, 0, 0, 1)`
    _, err = m.db.Exec(query, fields[1], fields[0], fields[3], price)
    if err != nil {
        m.status = fmt.Sprintf("Errore: %v", err)
        return m, nil
    }
    
    m.products = loadProducts(m.db)
    m.currentApp = StateProducts
    m.formData = ""
    m.status = "Prodotto salvato con successo"
    return m, nil
}

// View renders the UI
func (m MainModel) View() string {
    if m.err != nil {
        return lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FF0000")).
            Render(fmt.Sprintf("Errore: %v", m.err))
    }
    
    // Header
    header := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#00FFFF")).
        Padding(1, 2).
        Render("🏢 GAzie TUI - Gestione Aziendale")
    
    // Content
    var content string
    switch m.currentApp {
    case StateMainMenu:
        content = m.renderMainMenu()
    case StateCustomers:
        content = m.renderCustomers()
    case StateProducts:
        content = m.renderProducts()
    case StateCustomerForm:
        content = m.renderCustomerForm()
    case StateProductForm:
        content = m.renderProductForm()
    }
    
    // Status bar
    statusBar := lipgloss.NewStyle().
        Background(lipgloss.Color("#333333")).
        Foreground(lipgloss.Color("#FFFFFF")).
        Padding(0, 1).
        Render(m.status + " | ↑↓: Naviga | INVIO: Seleziona | N: Nuovo | ESC: Indietro/Esci")
    
    // Combine all parts
    return lipgloss.JoinVertical(
        lipgloss.Top,
        header,
        content,
        statusBar,
    )
}

// renderMainMenu renders the main menu
func (m MainModel) renderMainMenu() string {
    var menu strings.Builder
    menu.WriteString("Seleziona un'opzione:\n\n")
    
    for i, item := range m.menuItems {
        if i == m.selected {
            menu.WriteString("▶ ")
        } else {
            menu.WriteString("  ")
        }
        menu.WriteString(item)
        if i < len(m.menuItems)-1 {
            menu.WriteString("\n")
        }
    }
    
    return lipgloss.NewStyle().
        Width(m.width - 4).
        Height(m.height - 8).
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        Render(menu.String())
}

// renderCustomers renders the customers list
func (m MainModel) renderCustomers() string {
    var content strings.Builder
    content.WriteString("📋 GESTIONE CLIENTI\n\n")
    content.WriteString("Codice    Nome                    Città           Telefono\n")
    content.WriteString("--------  ----------------------  --------------  ------------\n")
    
    for i, customer := range m.customers {
        if i == m.selected {
            content.WriteString("▶ ")
        } else {
            content.WriteString("  ")
        }
        
        content.WriteString(fmt.Sprintf("%-8s  %-22s  %-14s  %-12s\n",
            customer.Code,
            truncateString(customer.Name, 22),
            truncateString(customer.City, 14),
            truncateString(customer.Phone, 12),
        ))
    }
    
    if len(m.customers) == 0 {
        content.WriteString("\nNessun cliente presente. Premi 'N' per aggiungerne uno.")
    }
    
    return lipgloss.NewStyle().
        Width(m.width - 4).
        Height(m.height - 8).
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        Render(content.String())
}

// renderProducts renders the products list
func (m MainModel) renderProducts() string {
    var content strings.Builder
    content.WriteString("📦 GESTIONE PRODOTTI\n\n")
    content.WriteString("Codice    Nome                    Prezzo    Giacenza\n")
    content.WriteString("--------  ----------------------  --------  --------\n")
    
    for i, product := range m.products {
        if i == m.selected {
            content.WriteString("▶ ")
        } else {
            content.WriteString("  ")
        }
        
        content.WriteString(fmt.Sprintf("%-8s  %-22s  €%-7.2f  %d\n",
            product.Code,
            truncateString(product.Name, 22),
            product.Price,
            product.Stock,
        ))
    }
    
    if len(m.products) == 0 {
        content.WriteString("\nNessun prodotto presente. Premi 'N' per aggiungerne uno.")
    }
    
    return lipgloss.NewStyle().
        Width(m.width - 4).
        Height(m.height - 8).
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        Render(content.String())
}

// renderCustomerForm renders the customer form
func (m MainModel) renderCustomerForm() string {
    title := "Nuovo Cliente"
    if m.formMode == "edit" {
        title = "Modifica Cliente"
    }
    
    content := fmt.Sprintf(`%s

Inserisci i dati del cliente separati da |:
Nome|Codice|Città|Telefono|Email

Dati attuali: %s

Esempio: Mario Rossi|C001|Roma|06-123456|mario@email.it

Premi INVIO per salvare, ESC per annullare`, title, m.formData)
    
    return lipgloss.NewStyle().
        Width(m.width - 4).
        Height(m.height - 8).
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        Render(content)
}

// renderProductForm renders the product form
func (m MainModel) renderProductForm() string {
    title := "Nuovo Prodotto"
    if m.formMode == "edit" {
        title = "Modifica Prodotto"
    }
    
    content := fmt.Sprintf(`%s

Inserisci i dati del prodotto separati da |:
Nome|Codice|Prezzo|Descrizione

Dati attuali: %s

Esempio: Prodotto示例|P001|15.50|Descrizione del prodotto

Premi INVIO per salvare, ESC per annullare`, title, m.formData)
    
    return lipgloss.NewStyle().
        Width(m.width - 4).
        Height(m.height - 8).
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        Render(content)
}

// Helper functions
func initDatabase(db *sql.DB) {
    // Create app data directory
    os.MkdirAll(AppDataDir, 0755)
    
    // Create tables
    queries := []string{
        `CREATE TABLE IF NOT EXISTS customers (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            code TEXT NOT NULL UNIQUE,
            name TEXT NOT NULL,
            city TEXT,
            phone TEXT,
            email TEXT,
            active BOOLEAN DEFAULT 1,
            created DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
        
        `CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            code TEXT NOT NULL UNIQUE,
            name TEXT NOT NULL,
            description TEXT,
            price REAL DEFAULT 0,
            stock INTEGER DEFAULT 0,
            min_stock INTEGER DEFAULT 0,
            active BOOLEAN DEFAULT 1,
            created DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
    }
    
    for _, query := range queries {
        db.Exec(query)
    }
    
    // Insert sample data if empty
    var count int
    db.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
    if count == 0 {
        sampleCustomers := []string{
            "ABC SRL|CLI001|Roma|06-123456|info@abc.it",
            "XYZ SPA|CLI002|Milano|02-789012|contatti@xyz.it",
            "DEF SRL|CLI003|Napoli|081-345678|office@def.it",
        }
        
        for _, customer := range sampleCustomers {
            fields := strings.Split(customer, "|")
            if len(fields) >= 5 {
                db.Exec("INSERT INTO customers (code, name, city, phone, email, active) VALUES (?, ?, ?, ?, ?, 1)",
                    fields[1], fields[0], fields[2], fields[3], fields[4])
            }
        }
    }
    
    db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
    if count == 0 {
        sampleProducts := []string{
            "Prodotto示例|PROD001|25.50|Questo è un prodotto示例",
            "Articolo示例|ART002|15.25|Altro articolo di qualità",
            "Merce示例|MERC003|99.99|Merce di lusso",
        }
        
        for _, product := range sampleProducts {
            fields := strings.Split(product, "|")
            if len(fields) >= 4 {
                price, _ := strconv.ParseFloat(fields[2], 64)
                db.Exec("INSERT INTO products (code, name, description, price, stock, min_stock, active) VALUES (?, ?, ?, ?, 10, 2, 1)",
                    fields[1], fields[0], fields[3], price)
            }
        }
    }
}

func loadCustomers(db *sql.DB) []Customer {
    var customers []Customer
    rows, err := db.Query("SELECT id, code, name, city, phone, email, active, created FROM customers ORDER BY name")
    if err != nil {
        return customers
    }
    defer rows.Close()
    
    for rows.Next() {
        var customer Customer
        var createdStr string
        rows.Scan(&customer.ID, &customer.Code, &customer.Name, &customer.City,
            &customer.Phone, &customer.Email, &customer.Active, &createdStr)
        
        customer.Created, _ = time.Parse("2006-01-02 15:04:05", createdStr)
        customers = append(customers, customer)
    }
    
    return customers
}

func loadProducts(db *sql.DB) []Product {
    var products []Product
    rows, err := db.Query("SELECT id, code, name, description, price, stock, min_stock, active, created FROM products ORDER BY name")
    if err != nil {
        return products
    }
    defer rows.Close()
    
    for rows.Next() {
        var product Product
        var createdStr string
        rows.Scan(&product.ID, &product.Code, &product.Name, &product.Description,
            &product.Price, &product.Stock, &product.MinStock, &product.Active, &createdStr)
        
        product.Created, _ = time.Parse("2006-01-02 15:04:05", createdStr)
        products = append(products, product)
    }
    
    return products
}

func truncateString(s string, length int) string {
    if len(s) <= length {
        return s
    }
    if length <= 3 {
        return s[:length]
    }
    return s[:length-3] + "..."
}

func main() {
    // Initialize database
    db, err := sql.Open("sqlite3", DBPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to open database: %v\n", err)
        os.Exit(1)
    }
    defer db.Close()
    
    // Enable foreign keys
    db.Exec("PRAGMA foreign_keys = ON")
    
    // Create main application model
    model := NewMainModel(db)
    
    // Setup styling
    lipgloss.SetHasDarkBackground(true)
    
    // Start the application
    if err := tea.NewProgram(model).Start(); err != nil {
        log.Fatal("Failed to start application:", err)
    }
}