# Models by レイヤードアーキテクチャ

## 概要
このパッケージには、POS QRシステムのデータモデルが含まれています。

## テスト実装状況

### ✅ テスト完了済みモデル
すべてのモデルについてテスト実装が完了しており、テスト結果はすべて成功しています。

| モデル  | テストファイル    | ステータス   |
| ------- | ----------------- | ------------ |
| Claims  | `claims_test.go`  | ✅ 完了・成功 |
| Manager | `manager_test.go` | ✅ 完了・成功 |
| Order   | `order_test.go`   | ✅ 完了・成功 |
| Seat    | `seat_test.go`    | ✅ 完了・成功 |
| Session | `session_test.go` | ✅ 完了・成功 |
| Status  | `status_test.go`  | ✅ 完了・成功 |
| Store   | `store_test.go`   | ✅ 完了・成功 |
| Utils   | `utils_test.go`   | ✅ 完了・成功 |

## チーム開発規範

### テスト実装について
- **必須**: 新しいモデルを追加する際は、対応するテストファイル（`*_test.go`）を必ず作成してください
- **品質保証**: すべてのテストが成功することを確認してからコミットしてください
- **カバレッジ**: モデルの主要な機能とエッジケースをカバーするテストを実装してください
- **命名規則**: テストファイルは `{model_name}_test.go` の形式で命名してください

### テスト実行方法
```bash
# 全テスト実行
go test ./...

# 特定のモデルのテスト実行
go test -v {model_name}_test.go {model_name}.go

# カバレッジ付きテスト実行
go test -cover ./...
```

### 継続的品質管理
- 新機能追加時は対応するテストケースも追加してください
- リファクタリング後は既存テストがすべて通ることを確認してください
- テストの保守性を考慮し、適切なテストデータとモックを使用してください

## Status移行ルール

### ステータス定義
注文の状態は以下のカテゴリに分類されます：

#### 初期状態と支払い
- `created` → `pending_payment`, `confirmed`, `cancelled`, `declined`
- `pending_payment` → `payment_received`, `payment_failed`, `cancelled`, `declined`
- `payment_received` → `pending_confirmation`, `cancelled`, `declined`
- `payment_failed` → `cancelled`

#### レストラン内部での処理
- `pending_confirmation` → `confirmed`, `cancelled`, `declined`
- `confirmed` → `preparing`, `on_hold`, `cancelled`
- `preparing` → `ready_for_pickup`, `ready_for_delivery`, `on_hold`, `cancelled`
- `on_hold` → `confirmed`, `preparing`, `ready_for_pickup`, `ready_for_delivery`, `cancelled`

#### 提供・引き渡し・配送
- `ready_for_pickup` → `picked_up`, `served`, `cancelled`
- `ready_for_delivery` → `out_for_delivery`, `served`, `cancelled`
- `out_for_delivery` → `delivered`, `delivery_attempt_failed`, `cancelled`
- `delivery_attempt_failed` → `out_for_delivery`, `delivered`, `cancelled`

#### 完了と例外処理
- `delivered`, `picked_up`, `served` → `completed`, `refunded`, `partially_refunded`
- `partially_refunded` → `refunded`, `completed`, `cancelled`
- `failed` → `cancelled`

### 最終状態
以下のステータスは最終状態であり、これ以上の遷移はできません：
- `completed`
- `cancelled`
- `declined`
- `refunded`
- `payment_failed`
- `failed`

### ビジネスルール
- **アイテム追加可能**: `created`, `pending_payment`, `pending_confirmation`, `confirmed`, `on_hold`
- **キャンセル可能**: `created`, `pending_payment`, `pending_confirmation`, `preparing`, `confirmed`, `on_hold`, `ready_for_pickup`, `ready_for_delivery`
- **提供準備完了**: `ready_for_pickup`, `ready_for_delivery`, `out_for_delivery`
- **（外部配達である場合）履行完了**: `delivered`, `picked_up`, `served`
- **調理中**: `confirmed`, `preparing`, `on_hold`

### 実装上の注意点
- 状態遷移は `CanTransitionTo(newStatus Status) bool` メソッドで検証してください
- 不正な遷移を防ぐため、必ず事前チェックを行ってください
- 最終状態からの遷移は原則として禁止されています
- `partially_refunded` は例外的に限定的な遷移が許可されています