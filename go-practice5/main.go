package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" 
)

type ProductDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

func main() {

	dsn := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/shop?sslmode=disable")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("open db:", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("ping db:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/products", productsHandler(db))

	addr := ":8081"
	log.Println("server is running on", addr)
	if err := http.ListenAndServe(addr, logRequest(mux)); err != nil {
		log.Fatal(err)
	}
}


func productsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		base := `
				SELECT p.id, p.name, c.name AS category, p.price
				FROM products p
				JOIN categories c ON c.id = p.category_id
				WHERE 1=1
				`
		conds := make([]string, 0, 4)
		args := make([]any, 0, 6)
		argPos := 1 

		if cat := strings.TrimSpace(q.Get("category")); cat != "" {
			conds = append(conds, "AND c.name = $"+itoa(argPos))
			args = append(args, cat)
			argPos++
		}

		if s := strings.TrimSpace(q.Get("min_price")); s != "" {
			if v, err := strconv.Atoi(s); err == nil {
				conds = append(conds, "AND p.price >= $"+itoa(argPos))
				args = append(args, v)
				argPos++
			}
		}

		if s := strings.TrimSpace(q.Get("max_price")); s != "" {
			if v, err := strconv.Atoi(s); err == nil {
				conds = append(conds, "AND p.price <= $"+itoa(argPos))
				args = append(args, v)
				argPos++
			}
		}

	
		sort := strings.ToLower(strings.TrimSpace(q.Get("sort")))
		orderBy := ""
		switch sort {
		case "price_asc":
			orderBy = " ORDER BY p.price ASC"
		case "price_desc":
			orderBy = " ORDER BY p.price DESC"
		default:
			orderBy = " ORDER BY p.id ASC"
		}


		limit := 50
		if s := q.Get("limit"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v > 0 && v <= 500 {
				limit = v
			}
		}
		offset := 0
		if s := q.Get("offset"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v >= 0 {
				offset = v
			}
		}

		sqlStr := base + " " + strings.Join(conds, " ") + orderBy +
			" LIMIT $" + itoa(argPos) + " OFFSET $" + itoa(argPos+1)
		args = append(args, limit, offset)

		
		start := time.Now()
		rows, err := db.QueryContext(r.Context(), sqlStr, args...)
		queryTime := time.Since(start)
		w.Header().Set("X-Query-Time", queryTime.String())
		if err != nil {
			http.Error(w, "db query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var out []ProductDTO
		for rows.Next() {
			var p ProductDTO
			if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price); err != nil {
				http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			out = append(out, p)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "rows error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.String(), time.Since(t0))
	})
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func itoa(i int) string { return strconv.Itoa(i) }
