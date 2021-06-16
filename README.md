🐯 API Server for Code Archives 🐯

【使用技術】

- クライアントサイド
  - 言語： TypeScript 4.1.5
  - ライブラリ： React 17.0.1
  - フレームワーク: Next.js 10.0.6
  - デプロイ： Vercel
  - コードハイライト: Prism.js
  - CSS: tailwindcss 2.0.3
- サーバーサイド
  - 言語： Go 1.15
  - データベース: CloudSQL(MySQL 8.0)
  - デプロイ： Google App Engine

【機能一覧】

- コードの新規登録・編集・削除(現在文字数の多いコードの保存が適切に行えません)
- 複数言語に対応したコードハイライト機能
- 検索機能
- ユーザー登録・認証機能
