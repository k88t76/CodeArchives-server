package models

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"time"
)

type Archive struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid"`
	Content   string `json:"content"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Language  string `json:"language"`
	CreatedAt string `json:"createdAt"`
}

func NewArchive(id int64, uuid string, content string, title string, author string, language string, createdAt string) *Archive {
	return &Archive{
		id,
		uuid,
		content,
		title,
		author,
		language,
		createdAt,
	}
}

func GetArchive(uuid string) *Archive {
	cmd := fmt.Sprintf("SELECT id, uuid, content, title, author, language, created_at FROM %s WHERE uuid = ?", tableNameArchives)
	row := db.QueryRow(cmd, uuid)
	var archive Archive
	err := row.Scan(&archive.ID, &archive.UUID, &archive.Content, &archive.Title, &archive.Author, &archive.Language, &archive.CreatedAt)
	if err != nil {
		return nil
	}
	archive.CreatedAt = strings.Split(archive.CreatedAt, "T")[0]
	return NewArchive(archive.ID, archive.UUID, archive.Content, archive.Title, archive.Author, archive.Language, archive.CreatedAt)
}

func GetArchivesByUser(userName string, limit int) ([]Archive, error) {
	var archives []Archive
	cmd := fmt.Sprintf(`SELECT id, uuid, content, title, author, language, created_at FROM %s WHERE author = ? ORDER BY created_at DESC LIMIT ?`, tableNameArchives)
	rows, err := db.Query(cmd, userName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var archive Archive
		rows.Scan(&archive.ID, &archive.UUID, &archive.Content, &archive.Title, &archive.Author, &archive.Language, &archive.CreatedAt)
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		archive.CreatedAt = strings.Split(archive.CreatedAt, "T")[0]
		archives = append(archives, archive)
	}
	return archives, nil

}

func GetMatchArchive(search string, userName string) ([]Archive, error) {
	var archives []Archive
	cmd := fmt.Sprintf(`SELECT id, uuid, content, title, author, language, created_at FROM %s WHERE author = ? AND ( content LIKE `, tableNameArchives)
	words := strings.Fields(search)
	len := len(words)
	if len == 1 {
		cmd += "'%" + words[0] + "%')"
	} else {
		for i := 0; i < len-1; i++ {
			cmd += "'%" + words[0] + "%' OR content LIKE "
		}
		cmd += "'%" + words[len-1] + "%')"
	}
	rows, err := db.Query(cmd, userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var archive Archive
		rows.Scan(&archive.ID, &archive.UUID, &archive.Content, &archive.Title, &archive.Author, &archive.Language, &archive.CreatedAt)
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		archive.CreatedAt = strings.Split(archive.CreatedAt, "T")[0]
		archives = append(archives, archive)
	}
	return archives, nil

}

func (a *Archive) Create() error {
	if a.Content == "" {
		return nil
	}
	cmd := fmt.Sprintf("INSERT INTO %s (uuid, content, title, author, language, created_at) VALUES (?, ?, ?, ?, ?, ?)", tableNameArchives)
	_, err := db.Exec(cmd, CreateUUID(), a.Content, a.Title, a.Author, a.Language, time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02T15:04:05+09:00"))
	if err != nil {
		return err
	}
	return nil
}

func (a *Archive) Update() error {
	cmd := fmt.Sprintf("UPDATE %s SET uuid = ?, content = ?, title = ?, language = ?, created_at = ? WHERE id = ?", tableNameArchives)
	_, err := db.Exec(cmd, CreateUUID(), a.Content, a.Title, a.Language, time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02T15:04:05+09:00"), a.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *Archive) Delete() error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE uuid = ?", tableNameArchives)
	_, err := db.Exec(cmd, a.UUID)
	if err != nil {
		return err
	}
	return nil
}

func CreateUUID() string {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return uuid
}

func CreateGuestArchives() {
	guestArchives := []Archive{
		{Content: `package controllers

		import (
			"encoding/json"
			"fmt"
			"log"
			"net/http"
			"os"
			"path"
		
			"github.com/k88t76/CodeArchives-server/models"
		)
		
		func StartWebServer() {
			http.HandleFunc("/archive/", get)
			http.HandleFunc("/archives", getAll)
			http.HandleFunc("/create", create)
			http.HandleFunc("/edit/", edit)
			http.HandleFunc("/delete/", delete)
			http.HandleFunc("/search/", search)
			http.HandleFunc("/signin", signIn)
			http.HandleFunc("/signup", signUp)
			http.HandleFunc("/userbytoken", userByToken)
			http.HandleFunc("/testsignin", testSignIn)
		
			// [START setting_port]
			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
				log.Printf("Defaulting to port %s", port)
			}
		
			log.Printf("Listening on port %s", port)
			if err := http.ListenAndServe(":"+port, nil); err != nil {
				log.Fatal(err)
			}
			// [END setting_port]
		}
		
		func getAll(w http.ResponseWriter, r *http.Request) {
			setHeader(w)
			len := r.ContentLength
			body := make([]byte, len)
			r.Body.Read(body)
			var token string
			json.Unmarshal(body, &token)
			name, _ := models.GetUserNameByToken(token)
			archives, _ := models.GetArchivesByUser(name, 1000)
			output, err := json.MarshalIndent(&archives, "", "\t\t")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(output)
			return
		}`,
			Title:    "server.go",
			Author:   "guest-user",
			Language: "go",
		},
		{Content: `async fn maybe_fail_hello(setting: bool) -> Result<String, dyn std::error::Error> {
			if setting {
				Ok("hello".to_string())
			} else {
				Err("fail")
			}
		}
		
		#[tokio::main]
		async fn print_async_hello(setting: bool) {
			let greeting: String = maybe_fail_hello(setting).await?;
			println!("{}", greeting);
		}`,
			Title:    "async.rs",
			Author:   "guest-user",
			Language: "rust",
		},
		{Content: `SELECT 
		CONVERT(int, SUBSTRING(MIN(jis_code), 1, 2)) AS id
	  , pref
	  , COUNT(postal7) AS cnt 
	FROM POSTAL_CODE 
	WHERE pref IN ('東京都', '神奈川県', '埼玉県', '千葉県')
	GROUP BY pref
	UNION ALL
	SELECT 
		NULL AS id 
	  , '合計' AS pref
	  , SUM(cnt) AS cnt
	FROM (
	SELECT 
		CONVERT(int, SUBSTRING(MIN(jis_code), 1, 2)) AS id
	  , pref
	  , COUNT(postal7) AS cnt
	FROM POSTAL_CODE 
	WHERE pref IN ('東京都', '神奈川県', '埼玉県', '千葉県')
	GROUP BY pref
	) X
	ORDER BY id ASC 
		`,
			Title:    "SQL",
			Author:   "guest-user",
			Language: "sql",
		},
		{Content: `export async function getServerSideProps({res, req}) {
			const token = req.cookies.token || ''
			const data = await fetch('http://localhost:8080/archives', {
					  method: 'POST',
					  mode: 'cors',
					  headers: {'Content-Type': 'application/json'},
					  body: JSON.stringify(token),
					});
			console.log(data)
			const posts = await data.json()
			
			return {
			  props: {
				posts,
				token,
			  },
			}
		  }
		  
		  export default function Home ({posts, token})  {
		  
		  const [archives, setArchives] = useState(posts)
		  
		  
		  const [response, setResponse] = useState({
			type: '',
			message: '',
		  });
		  
		  const { query } = useRouter()
		  
		  useEffect(() => {
			setResponse({message: query.response})
			Prism.highlightAll()
		  }, []);
		  
		  
		  if(token === ''){
			return (
			  <Layout home>
				<HeaderUnLogin />
				<div className="content">
				<p>{response.message}</p>
				<div className='Signin'>
				<Link href='/signin'>
				  <a href="#" className="btn-signin">Sign in</a>
				</Link>
				<div className="or">or</div>
				<Link href='/test-signin'>
				  <a href="#" className="btn-signin test">Test Sign in</a>
				</Link>
				</div>
				</div>
			  </Layout>
			  
			)
		  } else {
			return(
		  
			  <Layout home>
				<header className={styles.header}>
				<HeaderLogin/>
			  </header>
			  <div className='content'>
		   <div>
			  <p>{response.message}</p>
		  
			</div>
		  
		  <section>
		  <h2 className={utilStyles.headingLg}>Archives</h2>
		  
		  <Search  setArchives={setArchives} token={token} archives={archives}/>
		  
		  {archives && (  
		  <Archives archives={archives} />
		  )}
		  
		  { !archives && (
			<div>Empty</div>
		  )}
		  </section>
		  </div>
		  </Layout>
			)
		  } 
		  }`,
			Title:    "index.js",
			Author:   "guest-user",
			Language: "jsx",
		},
		{Content: `FROM golang:1.12.0-alpine3.9

		WORKDIR /go/src/app
		
		ENV GO111MODULE=on
		
		COPY go.mod go.sum ./
		
		RUN go mod download
		
		COPY . .
		
		RUN go build -o main . 
		
		EXPOSE 8080
		
		CMD ["/go/src/app"]
		`,
			Title:    "Dockerfile",
			Author:   "guest-user",
			Language: "docker",
		},
		{Content: `package main

		import (
			"bufio"
			"fmt"
			"os"
			"strconv"
			"strings"
		)
		
		var sc = bufio.NewScanner(os.Stdin)
		
		func Scan() string {
			sc.Scan()
			return sc.Text()
		}
		func rScan() []rune {
			return []rune(Scan())
		}
		func iScan() int {
			n, _ := strconv.Atoi(Scan())
			return n
		}
		func fScan() float64 {
			n, _ := strconv.ParseFloat(Scan(), 64)
			return n
		}
		func SScan(n int) []string {
			a := make([]string, n)
			for i := 0; i < n; i++ {
				a[i] = Scan()
			}
			return a
		}
		func iSScan(n int) []int {
			a := make([]int, n)
			for i := 0; i < n; i++ {
				a[i] = iScan()
			}
			return a
		}
		func atoi(s string) int {
			n, _ := strconv.Atoi(s)
			return n
		}
		func abs(x int) int {
			if x < 0 {
				x = -x
			}
			return x
		}
		func mod(x, d int) int {
			x %= d
			if x < 0 {
				x += d
			}
			return x
		}
		func max(a ...int) int {
			x := -int(1e+18)
			for i := 0; i < len(a); i++ {
				if x < a[i] {
					x = a[i]
				}
			}
			return x
		}
		func min(a ...int) int {
			x := int(1e+18)
			for i := 0; i < len(a); i++ {
				if x > a[i] {
					x = a[i]
				}
			}
			return x
		}
		func sum(a []int) int {
			x := 0
			for _, v := range a {
				x += v
			}
			return x
		}
		func fSum(a []float64) float64 {
			x := 0.
			for _, v := range a {
				x += v
			}
			return x
		}
		func bPrint(f bool, x string, y string) {
			if f {
				fmt.Println(x)
			} else {
				fmt.Println(y)
			}
		}
		func iSSPrint(x []int) {
			fmt.Println(strings.Trim(fmt.Sprint(x), "[]"))
		}
		
		var P1 int = 1000000007
		var P2 int = 998244353
		
		func main() {
			buf := make([]byte, 0)
			sc.Buffer(buf, P1)
			sc.Split(bufio.ScanWords)
			n := iScan()
			s := SScan(n)
			t, f := make([]int, n+1), make([]int, n+1)
			t[0], f[0] = 1, 1
			for i := 1; i <= n; i++ {
				if s[i-1] == "AND" {
					t[i] = t[i-1]
					f[i] = t[i-1] + f[i-1]*2
				} else {
					t[i] = t[i-1]*2 + f[i-1]
					f[i] = f[i-1]
				}
			}
			fmt.Println(t[n])
		}`,
			Title:    "algorithm",
			Author:   "guest-user",
			Language: "go",
		},
		{Content: `import numpy as np
		import math
		 
		def is_prime(n):
			if n % 2 == 0 and n > 2: 
				return False
			return all(n % i for i in range(3, int(math.sqrt(n)) + 1, 2))
		 
		arr = np.arange(2, 21)
		vec = np.vectorize(is_prime)
		print(vec(arr))
		print(arr[vec(arr)])`,
			Title:    "prime.py",
			Author:   "guest-user",
			Language: "python",
		},
		{Content: `for i in 1..100
		if i % 15 == 0
			print "FizzBuzz\s"
		elsif i % 5 == 0
			print "Buzz\s"
		elsif i % 3 == 0
			print "Fizz\s"
		else
			print i , "\s"
		end
		end`,
			Title:    "fizzbuzz.rb",
			Author:   "guest-user",
			Language: "ruby",
		},
		{Content: `<!DOCTYPE html>
		<html lang="ja">
		 <head>
		 <meta charset="utf-8">
		 <title>サイトタイトル</title>
		 <meta name="description" content="ディスクリプションを入力">
		 <meta name="viewport" content="width=device-width, initial-scale=1.0">
		 <link rel="stylesheet" href="style.css">
		 <!-- [if lt IE 9] -->
		 <script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
		 <script src="http://css3-mediaqueries-js.googlecode.com/svn/trunk/css3-mediaqueries.js"></script>
		 <!-- [endif] -->
		 <script src="main.js"></script>
		 </head>
		 <body>
		 <!----- header----->
		 <header>ヘッダー</header>
		 <nav>ナビ</nav>
		 <!----- /header ----->
		 
		 <!----- main ----->
		 <article>
		 <h1>タイトル</h1>
		 <section>
		 <h2>見出し２</h2>
		 <p>コンテンツの内容</p>
		 </section>
		 </article>
		 <!----- /main ----->
		 
		 <!----- footer ----->
		 <footer>フッター</footer>
		 <!----- /footer ----->
		 </body>
		</html>`,
			Title:    "template.html",
			Author:   "guest-user",
			Language: "html",
		},
		{Content: `package main

		import "fmt"
		
		func main() {
			fmt.Println("Hello, World!")
		}`,
			Title:    "Hello World",
			Author:   "guest-user",
			Language: "go",
		},
	}
	cmd := fmt.Sprintf("DELETE FROM %s WHERE author = ?", tableNameArchives)
	db.Exec(cmd, "guest-user")
	for _, a := range guestArchives {
		cmd = fmt.Sprintf("INSERT INTO %s (uuid, content, title, author, language, created_at) VALUES (?, ?, ?, ?, ?, ?)", tableNameArchives)
		db.Exec(cmd, CreateUUID(), a.Content, a.Title, a.Author, a.Language, time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02T15:04:05+09:00"))
	}
}
