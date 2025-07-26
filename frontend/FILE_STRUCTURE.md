# フロントエンド ファイル構造定義書

## プロジェクト全体構造

```
frontend/
├── src/                           # ソースコード
├── public/                        # 静的ファイル
├── docs/                          # ドキュメント
├── tests/                         # テストファイル
├── .env.local                     # 環境変数（開発用）
├── .env.example                   # 環境変数テンプレート
├── next.config.js                 # Next.js設定
├── tailwind.config.js             # Tailwind設定
├── tsconfig.json                  # TypeScript設定
├── package.json                   # 依存関係
├── DEVELOPMENT_ROLES.md           # 開発役務定義
├── FILE_STRUCTURE.md              # ファイル構造定義（本ファイル）
└── README.md                      # プロジェクト概要
```

## src/ ディレクトリ詳細構造

### 1. app/ - Next.js App Router ページ

```
src/app/
├── (auth)/                        # 認証が必要なページグループ
│   ├── admin/                     # 店舗発行管理者用ページ
│   │   ├── dashboard/
│   │   │   ├── page.tsx          # 管理者ダッシュボード
│   │   │   └── loading.tsx       # ローディング画面
│   │   ├── stores/
│   │   │   ├── page.tsx          # 店舗一覧
│   │   │   ├── [id]/
│   │   │   │   ├── page.tsx      # 店舗詳細
│   │   │   │   └── edit/
│   │   │   │       └── page.tsx  # 店舗編集
│   │   │   ├── create/
│   │   │   │   └── page.tsx      # 店舗作成
│   │   │   └── loading.tsx
│   │   ├── managers/
│   │   │   ├── page.tsx          # 管理者一覧
│   │   │   ├── [id]/
│   │   │   │   └── page.tsx      # 管理者詳細
│   │   │   ├── create/
│   │   │   │   └── page.tsx      # 管理者作成
│   │   │   └── loading.tsx
│   │   ├── layout.tsx            # 管理者レイアウト
│   │   └── loading.tsx
│   ├── store/                     # 店舗管理者用ページ
│   │   ├── dashboard/
│   │   │   ├── page.tsx          # 店舗ダッシュボード
│   │   │   └── loading.tsx
│   │   ├── seats/
│   │   │   ├── page.tsx          # 座席管理
│   │   │   ├── [id]/
│   │   │   │   ├── page.tsx      # 座席詳細
│   │   │   │   └── qr/
│   │   │   │       └── page.tsx  # QRコード表示
│   │   │   ├── create/
│   │   │   │   └── page.tsx      # 座席作成
│   │   │   └── loading.tsx
│   │   ├── orders/
│   │   │   ├── page.tsx          # 注文管理
│   │   │   ├── [id]/
│   │   │   │   └── page.tsx      # 注文詳細
│   │   │   └── loading.tsx
│   │   ├── menu/
│   │   │   ├── page.tsx          # メニュー管理
│   │   │   ├── [id]/
│   │   │   │   ├── page.tsx      # メニュー詳細
│   │   │   │   └── edit/
│   │   │   │       └── page.tsx  # メニュー編集
│   │   │   ├── create/
│   │   │   │   └── page.tsx      # メニュー作成
│   │   │   ├── categories/
│   │   │   │   └── page.tsx      # カテゴリ管理
│   │   │   └── loading.tsx
│   │   ├── layout.tsx            # 店舗レイアウト
│   │   └── loading.tsx
│   └── layout.tsx                # 認証レイアウト
├── auth/                          # 認証ページ
│   ├── admin-login/
│   │   └── page.tsx              # 管理者ログイン
│   ├── store-login/
│   │   └── page.tsx              # 店舗ログイン
│   └── layout.tsx                # 認証レイアウト
├── order/                         # 顧客注文ページ
│   ├── [sessionId]/
│   │   ├── page.tsx              # 注文画面
│   │   ├── menu/
│   │   │   └── page.tsx          # メニュー表示
│   │   ├── cart/
│   │   │   └── page.tsx          # カート画面
│   │   ├── history/
│   │   │   └── page.tsx          # 注文履歴
│   │   └── loading.tsx
│   └── layout.tsx                # 注文レイアウト
├── api/                           # API Routes（必要に応じて）
├── globals.css                    # グローバルスタイル
├── layout.tsx                     # ルートレイアウト
├── page.tsx                       # ホームページ
├── loading.tsx                    # グローバルローディング
├── error.tsx                      # グローバルエラー
└── not-found.tsx                  # 404ページ
```

### 2. components/ - 再利用可能コンポーネント

```
src/components/
├── ui/                            # 基本UIコンポーネント
│   ├── Button/
│   │   ├── index.tsx             # Button コンポーネント
│   │   ├── Button.types.ts       # 型定義
│   │   └── Button.module.css     # スタイル（必要に応じて）
│   ├── Card/
│   │   ├── index.tsx
│   │   └── Card.types.ts
│   ├── Modal/
│   │   ├── index.tsx
│   │   └── Modal.types.ts
│   ├── Table/
│   │   ├── index.tsx
│   │   └── Table.types.ts
│   ├── Form/
│   │   ├── index.tsx
│   │   └── Form.types.ts
│   └── index.ts                  # エクスポート集約
├── layouts/                       # レイアウトコンポーネント
│   ├── AdminLayout/
│   │   ├── index.tsx             # 管理者レイアウト
│   │   ├── Sidebar.tsx           # サイドバー
│   │   ├── Header.tsx            # ヘッダー
│   │   └── AdminLayout.types.ts
│   ├── StoreLayout/
│   │   ├── index.tsx             # 店舗レイアウト
│   │   ├── Sidebar.tsx
│   │   ├── Header.tsx
│   │   └── StoreLayout.types.ts
│   ├── AuthLayout/
│   │   ├── index.tsx             # 認証レイアウト
│   │   └── AuthLayout.types.ts
│   ├── OrderLayout/
│   │   ├── index.tsx             # 注文レイアウト
│   │   └── OrderLayout.types.ts
│   └── index.ts
├── forms/                         # フォームコンポーネント
│   ├── LoginForm/
│   │   ├── index.tsx             # ログインフォーム
│   │   ├── AdminLoginForm.tsx    # 管理者ログイン
│   │   ├── StoreLoginForm.tsx    # 店舗ログイン
│   │   └── LoginForm.types.ts
│   ├── StoreForm/
│   │   ├── index.tsx             # 店舗フォーム
│   │   ├── CreateStoreForm.tsx   # 店舗作成
│   │   ├── EditStoreForm.tsx     # 店舗編集
│   │   └── StoreForm.types.ts
│   ├── SeatForm/
│   │   ├── index.tsx
│   │   └── SeatForm.types.ts
│   ├── MenuForm/
│   │   ├── index.tsx
│   │   └── MenuForm.types.ts
│   └── index.ts
├── features/                      # 機能別コンポーネント
│   ├── auth/
│   │   ├── AuthGuard/
│   │   │   ├── index.tsx         # 認証ガード
│   │   │   └── AuthGuard.types.ts
│   │   ├── LoginButton/
│   │   │   └── index.tsx
│   │   └── LogoutButton/
│   │       └── index.tsx
│   ├── stores/
│   │   ├── StoreTable/
│   │   │   ├── index.tsx         # 店舗テーブル
│   │   │   └── StoreTable.types.ts
│   │   ├── StoreCard/
│   │   │   ├── index.tsx         # 店舗カード
│   │   │   └── StoreCard.types.ts
│   │   ├── StoreDetail/
│   │   │   ├── index.tsx         # 店舗詳細
│   │   │   └── StoreDetail.types.ts
│   │   └── StoreStats/
│   │       ├── index.tsx         # 店舗統計
│   │       └── StoreStats.types.ts
│   ├── seats/
│   │   ├── SeatGrid/
│   │   │   ├── index.tsx         # 座席グリッド
│   │   │   └── SeatGrid.types.ts
│   │   ├── SeatCard/
│   │   │   ├── index.tsx         # 座席カード
│   │   │   └── SeatCard.types.ts
│   │   ├── QRCodeDisplay/
│   │   │   ├── index.tsx         # QRコード表示
│   │   │   └── QRCodeDisplay.types.ts
│   │   └── SeatStatus/
│   │       ├── index.tsx         # 座席状態
│   │       └── SeatStatus.types.ts
│   ├── orders/
│   │   ├── OrderTable/
│   │   │   ├── index.tsx         # 注文テーブル
│   │   │   └── OrderTable.types.ts
│   │   ├── OrderCard/
│   │   │   ├── index.tsx         # 注文カード
│   │   │   └── OrderCard.types.ts
│   │   ├── OrderDetail/
│   │   │   ├── index.tsx         # 注文詳細
│   │   │   └── OrderDetail.types.ts
│   │   ├── OrderStatusBadge/
│   │   │   ├── index.tsx         # 注文ステータス
│   │   │   └── OrderStatusBadge.types.ts
│   │   ├── OrderCart/
│   │   │   ├── index.tsx         # 注文カート
│   │   │   └── OrderCart.types.ts
│   │   └── OrderHistory/
│   │       ├── index.tsx         # 注文履歴
│   │       └── OrderHistory.types.ts
│   ├── menu/
│   │   ├── MenuGrid/
│   │   │   ├── index.tsx         # メニューグリッド
│   │   │   └── MenuGrid.types.ts
│   │   ├── MenuCard/
│   │   │   ├── index.tsx         # メニューカード
│   │   │   └── MenuCard.types.ts
│   │   ├── MenuTable/
│   │   │   ├── index.tsx         # メニューテーブル
│   │   │   └── MenuTable.types.ts
│   │   └── CategoryManager/
│   │       ├── index.tsx         # カテゴリ管理
│   │       └── CategoryManager.types.ts
│   ├── dashboard/
│   │   ├── StatCard/
│   │   │   ├── index.tsx         # 統計カード
│   │   │   └── StatCard.types.ts
│   │   ├── Chart/
│   │   │   ├── index.tsx         # チャート
│   │   │   └── Chart.types.ts
│   │   ├── ActivityList/
│   │   │   ├── index.tsx         # アクティビティリスト
│   │   │   └── ActivityList.types.ts
│   │   └── RevenueChart/
│   │       ├── index.tsx         # 売上チャート
│   │       └── RevenueChart.types.ts
│   └── index.ts
└── index.ts                       # 全コンポーネントエクスポート
```

### 3. hooks/ - カスタムフック

```
src/hooks/
├── auth/
│   ├── useAuth.ts                # 認証フック
│   ├── useLogin.ts               # ログインフック
│   └── useLogout.ts              # ログアウトフック
├── api/
│   ├── useStores.ts              # 店舗API フック
│   ├── useManagers.ts            # 管理者API フック
│   ├── useSeats.ts               # 座席API フック
│   ├── useOrders.ts              # 注文API フック
│   ├── useMenu.ts                # メニューAPI フック
│   └── useSessions.ts            # セッションAPI フック
├── ui/
│   ├── useModal.ts               # モーダルフック
│   ├── useTable.ts               # テーブルフック
│   ├── useForm.ts                # フォームフック
│   └── useNotification.ts       # 通知フック
├── utils/
│   ├── useLocalStorage.ts        # ローカルストレージフック
│   ├── useDebounce.ts            # デバウンスフック
│   └── useMediaQuery.ts          # メディアクエリフック
└── index.ts                      # エクスポート集約
```

### 4. lib/ - ユーティリティ・設定

```
src/lib/
├── api/
│   ├── client.ts                 # API クライアント設定
│   ├── endpoints.ts              # API エンドポイント定義
│   ├── interceptors.ts           # リクエスト/レスポンス インターセプター
│   └── types.ts                  # API 型定義
├── auth/
│   ├── jwt.ts                    # JWT ユーティリティ
│   ├── storage.ts                # 認証情報ストレージ
│   ├── guards.ts                 # 認証ガード
│   └── types.ts                  # 認証型定義
├── utils/
│   ├── format.ts                 # フォーマット関数
│   ├── validation.ts             # バリデーション関数
│   ├── constants.ts              # 定数定義
│   ├── helpers.ts                # ヘルパー関数
│   └── types.ts                  # 共通型定義
├── config/
│   ├── env.ts                    # 環境変数設定
│   ├── theme.ts                  # テーマ設定
│   └── routes.ts                 # ルート定義
└── index.ts                      # エクスポート集約
```

### 5. store/ - 状態管理

```
src/store/
├── auth/
│   ├── authStore.ts              # 認証状態ストア
│   ├── authTypes.ts              # 認証状態型定義
│   └── authActions.ts            # 認証アクション
├── ui/
│   ├── uiStore.ts                # UI状態ストア
│   ├── uiTypes.ts                # UI状態型定義
│   └── uiActions.ts              # UIアクション
├── data/
│   ├── storeStore.ts             # 店舗データストア
│   ├── orderStore.ts             # 注文データストア
│   └── menuStore.ts              # メニューデータストア
├── providers/
│   ├── QueryProvider.tsx        # React Query プロバイダー
│   ├── AuthProvider.tsx         # 認証プロバイダー
│   └── ThemeProvider.tsx        # テーマプロバイダー
├── index.ts                      # ストア設定
└── types.ts                      # 状態管理型定義
```

### 6. types/ - TypeScript型定義

```
src/types/
├── api/
│   ├── auth.ts                   # 認証API型
│   ├── stores.ts                 # 店舗API型
│   ├── seats.ts                  # 座席API型
│   ├── orders.ts                 # 注文API型
│   ├── menu.ts                   # メニューAPI型
│   └── common.ts                 # 共通API型
├── models/
│   ├── User.ts                   # ユーザーモデル
│   ├── Store.ts                  # 店舗モデル
│   ├── Seat.ts                   # 座席モデル
│   ├── Order.ts                  # 注文モデル
│   ├── Menu.ts                   # メニューモデル
│   └── Session.ts                # セッションモデル
├── ui/
│   ├── components.ts             # コンポーネント型
│   ├── forms.ts                  # フォーム型
│   └── layouts.ts                # レイアウト型
├── utils/
│   ├── helpers.ts                # ヘルパー型
│   └── constants.ts              # 定数型
├── global.d.ts                   # グローバル型定義
└── index.ts                      # 型エクスポート集約
```

## 静的ファイル構造

### public/ ディレクトリ

```
public/
├── images/
│   ├── logo/
│   │   ├── logo.svg              # ロゴ
│   │   └── favicon.ico           # ファビコン
│   ├── icons/
│   │   ├── menu-icons/           # メニューアイコン
│   │   └── ui-icons/             # UIアイコン
│   └── placeholders/
│       ├── menu-placeholder.jpg  # メニュー画像プレースホルダー
│       └── store-placeholder.jpg # 店舗画像プレースホルダー
├── fonts/                        # カスタムフォント（必要に応じて）
└── manifest.json                 # PWA マニフェスト
```

## 設定ファイル詳細

### package.json 依存関係

```json
{
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.0.0",
    "react-dom": "^18.0.0",
    "antd": "^5.0.0",
    "@ant-design/nextjs-registry": "^1.0.0",
    "@tanstack/react-query": "^5.0.0",
    "axios": "^1.6.0",
    "zustand": "^4.4.0",
    "js-cookie": "^3.0.0",
    "react-hook-form": "^7.47.0",
    "@hookform/resolvers": "^3.3.0",
    "zod": "^3.22.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/react": "^18.0.0",
    "@types/react-dom": "^18.0.0",
    "@types/js-cookie": "^3.0.0",
    "typescript": "^5.0.0",
    "eslint": "^8.0.0",
    "eslint-config-next": "^14.0.0",
    "prettier": "^3.0.0",
    "@typescript-eslint/eslint-plugin": "^6.0.0",
    "@typescript-eslint/parser": "^6.0.0"
  }
}
```

### tsconfig.json 設定

```json
{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "es6"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"],
      "@/components/*": ["./src/components/*"],
      "@/hooks/*": ["./src/hooks/*"],
      "@/lib/*": ["./src/lib/*"],
      "@/store/*": ["./src/store/*"],
      "@/types/*": ["./src/types/*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

## ファイル命名規則

### コンポーネントファイル

- **ディレクトリ名**: PascalCase (`StoreTable/`)
- **ファイル名**: `index.tsx` (メインコンポーネント)
- **型定義**: `ComponentName.types.ts`
- **スタイル**: `ComponentName.module.css` (必要に応じて)

### フック

- **ファイル名**: camelCase (`useAuth.ts`)
- **関数名**: camelCase (`useAuth`)

### ユーティリティ

- **ファイル名**: camelCase (`formatDate.ts`)
- **関数名**: camelCase (`formatDate`)

### 型定義

- **ファイル名**: PascalCase (`User.ts`)
- **型名**: PascalCase (`User`, `UserResponse`)
- **インターフェース名**: PascalCase (`IUser`)

### 定数

- **ファイル名**: camelCase (`apiEndpoints.ts`)
- **定数名**: UPPER_SNAKE_CASE (`API_BASE_URL`)

## エクスポート戦略

### バレルエクスポート

各ディレクトリに `index.ts` を配置し、関連するモジュールを集約エクスポート

```typescript
// components/index.ts
export { default as Button } from './ui/Button';
export { default as StoreTable } from './features/stores/StoreTable';
export { default as AdminLayout } from './layouts/AdminLayout';

// hooks/index.ts
export { useAuth } from './auth/useAuth';
export { useStores } from './api/useStores';

// types/index.ts
export type { User } from './models/User';
export type { Store } from './models/Store';
```

この構造により、開発効率と保守性を両立した、スケーラブルなフロントエンドアプリケーションを構築できます。
