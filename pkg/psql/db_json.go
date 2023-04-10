package psql

type StDBJson struct {
    Host     string    `json:"host"`
    User     string    `json:"user"`
    Password string    `json:"password"`
    DBName   string    `json:"dbname"`
    Tables   []StTable `json:"tables"`
}

type StTable struct {
    TableName string `json:"name"`
    TableFile string `json:"file"`
    CreateExp string `json:"create"`
}
