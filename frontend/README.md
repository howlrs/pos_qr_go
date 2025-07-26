# POS QR System - Frontend

Next.js フロントエンドアプリケーション for POS QRシステム

## 🚀 技術スタック

- **Framework**: Next.js 15.4.4 (App Router)
- **Language**: TypeScript 5.x
- **UI Library**: Ant Design 5.x
- **Styling**: Tailwind CSS 4.x
- **State Management**: Zustand + React Query
- **HTTP Client**: Axios
- **Form Handling**: React Hook Form + Zod
- **Authentication**: JWT (Cookie-based)

## 📁 プロジェクト構造

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router pages
│   ├── components/             # 再利用可能コンポーネント
│   ├── hooks/                  # カスタムフック
│   ├── lib/                    # ユーティリティ・設定
│   ├── store/                  # 状態管理
│   └── types/                  # TypeScript型定義
├── public/                     # 静的ファイル
├── DEVELOPMENT_ROLES.md        # 開発役務定義
├── FILE_STRUCTURE.md           # ファイル構造定義
└── DIRECTION.md                # 開発ディレクション
```

## 🛠️ 開発環境セットアップ

### 前提条件

- Node.js 18.x以上
- npm

### インストール

```bash
# 依存関係のインストール
npm install

# 環境変数の設定
cp .env.example .env.local
# .env.local を編集して適切な値を設定
```

### 開発サーバー起動

```bash
# 開発サーバー起動
npm run dev

# ブラウザで http://localhost:3000 を開く
```

### その他のコマンド

```bash
# プロダクションビルド
npm run build

# プロダクションサーバー起動
npm run start

# リンター実行
npm run lint
```

## 🔧 環境変数

`.env.local` ファイルで以下の環境変数を設定してください：

```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=POS QR System
NEXT_PUBLIC_ENVIRONMENT=development

# JWT Configuration
JWT_SECRET=your-jwt-secret-key-here
```

## 📚 開発ドキュメント

- **[DEVELOPMENT_ROLES.md](./DEVELOPMENT_ROLES.md)** - 開発役務と責務定義
- **[FILE_STRUCTURE.md](./FILE_STRUCTURE.md)** - 詳細なファイル構造とコーディング規約
- **[DIRECTION.md](./DIRECTION.md)** - 開発ディレクションと概要

## 🎯 開発フロー

### Phase 1: 基盤構築

- [x] Next.js + TypeScript セットアップ
- [x] Ant Design 統合
- [x] 基本ディレクトリ構造作成
- [ ] 認証システム実装
- [ ] 基本レイアウト作成

### Phase 2: 管理者機能

- [ ] 店舗発行管理者画面
- [ ] 店舗管理機能
- [ ] 管理者管理機能

### Phase 3: 店舗管理機能

- [ ] 店舗ダッシュボード
- [ ] 座席管理・QRコード生成
- [ ] 注文管理画面

### Phase 4: 顧客注文機能

- [ ] 注文画面実装
- [ ] セッション管理
- [ ] 注文フロー完成

### Phase 5: 統合・最適化

- [ ] 全機能統合テスト
- [ ] パフォーマンス最適化
- [ ] デプロイ準備

## 🔗 関連リポジトリ

- **Backend API**: `../` (Go + Echo + Firestore)

## 📝 開発ガイドライン

1. **コンポーネント作成**: `FILE_STRUCTURE.md` の命名規則に従う
2. **型定義**: 厳格な TypeScript 型定義を使用
3. **状態管理**: React Query でサーバー状態、Zustand でクライアント状態
4. **スタイリング**: Ant Design コンポーネントを優先使用
5. **API連携**: Axios + React Query パターンを使用

## 🐛 トラブルシューティング

### よくある問題

1. **ビルドエラー**: `npm run build` でエラーが発生する場合
   - `node_modules` を削除して `npm install` を再実行
   - TypeScript エラーを確認

2. **API接続エラー**: バックエンドAPIに接続できない場合
   - `.env.local` の `NEXT_PUBLIC_API_URL` を確認
   - バックエンドサーバーが起動しているか確認

3. **Ant Design スタイルが適用されない**:
   - `next.config.ts` の設定を確認
   - ブラウザのキャッシュをクリア

## 📞 サポート

開発に関する質問や問題がある場合は、プロジェクトドキュメントを参照するか、開発チームにお問い合わせください。
