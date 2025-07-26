# Next.js フロントエンド開発 概要資料

## プロジェクト概要

POS QRシステムのフロントエンド開発

- **管理者用画面**: 店舗発行管理者・店舗管理者向けの管理インターフェース
- **顧客用画面**: QRコードアクセス後の注文インターフェース

## 技術スタック

### 基本構成

- **Framework**: Next.js (App Router, TypeScript)
- **UI Library**: Ant Design 5.x
- **Styling**: Ant Design標準カラーパレット
- **State Management**: React Query (TanStack Query) + Zustand
- **HTTP Client**: Axios
- **Form Handling**: Ant Design Form + React Hook Form
- **Authentication**: JWT (Cookie-based)

### 開発環境

- **Node.js**: 18.x以上
- **Package Manager**: npm
- **Linting**: ESLint + Prettier
- **Type Checking**: TypeScript 5.x

## 関連ドキュメント

### 📋 開発役務定義

詳細な開発役務と責務については以下を参照：

- **[DEVELOPMENT_ROLES.md](./DEVELOPMENT_ROLES.md)** - 開発役務一覧と責務定義

### 📁 ファイル構造定義

プロジェクトのファイル構造と命名規則については以下を参照：

- **[FILE_STRUCTURE.md](./FILE_STRUCTURE.md)** - 詳細なディレクトリ構造とファイル配置

## 開発マイルストーン

### Phase 1: 基盤構築 (Week 1)

- [ ] Next.js + Ant Design セットアップ
- [ ] 認証システム実装
- [ ] 基本レイアウト作成

### Phase 2: 管理者機能 (Week 2-3)

- [ ] 店舗発行管理者画面
- [ ] 店舗管理機能
- [ ] 管理者管理機能

### Phase 3: 店舗管理機能 (Week 4-5)

- [ ] 店舗ダッシュボード
- [ ] 座席管理・QRコード生成
- [ ] 注文管理画面

### Phase 4: 顧客注文機能 (Week 6)

- [ ] 注文画面実装
- [ ] セッション管理
- [ ] 注文フロー完成

### Phase 5: 統合・最適化 (Week 7)

- [ ] 全機能統合テスト
- [ ] パフォーマンス最適化
- [ ] デプロイ準備

## 環境変数設定

```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=POS QR System
NEXT_PUBLIC_ENVIRONMENT=development

# .env.production
NEXT_PUBLIC_API_URL=https://your-api-domain.com
NEXT_PUBLIC_APP_NAME=POS QR System
NEXT_PUBLIC_ENVIRONMENT=production
```

## 品質基準

- **TypeScript**: 厳格な型チェック、型安全性100%
- **コードカバレッジ**: 80%以上
- **パフォーマンス**: Lighthouse スコア90以上
- **アクセシビリティ**: WCAG 2.1 AA準拠
- **レスポンシブ**: 全デバイス対応
- **ブラウザ対応**: Chrome, Firefox, Safari, Edge 最新版

## 開発開始手順

1. **役務確認**: `DEVELOPMENT_ROLES.md` で担当役務を確認
2. **構造理解**: `FILE_STRUCTURE.md` でファイル配置を理解
3. **環境構築**: 必要な依存関係をインストール
4. **開発開始**: Phase 1から順次実装

詳細な実装指針や技術仕様については、各専門ドキュメントを参照してください。
