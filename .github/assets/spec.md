# 混雑状況投票 (Crowd Vote) API 仕様書

## 概要

Crowd Vote は、ユーザーがさまざまな場所の「混雑具合」をリアルタイムで報告・確認できるサービスです。
このAPIは、場所の管理とユーザーからの投票を管理するためのバックエンド機能を提供します。
基本的にコンテナ上で動作するサーバレスなAPIを想定しており、コンテナ上のデータベースにデータを保存します。

## データモデル

### 場所 (Location)

- `group`: 場所のグループ (例: "tokyo", "osaka")
- `id`: 一意識別子
- `name`: 場所の名前 (例: "中央通りカフェ")
- `description`: 説明 (任意)

### 投票 (Vote)

- `group`: 場所のグループへの参照
- `location_id`: 場所への参照ID
- `client_address`: 投票を行ったユーザーのIPアドレス
- `level`: 報告された混雑度 (0-3)
  - `0`: ガラガラ (Empty)
  - `1`: 普通 (Normal)
  - `2`: 混雑 (Crowded)
  - `3`: 非常に混雑 (Very Crowded)
- `timestamp`: 報告日時

### CrowdRate

- `gruop`: 場所のグループへの参照
- `location_id`: 場所への参照ID
- `votes`: 指定された時間枠内の投票の加重平均を表します（float型）。
- `total`: 指定された時間枠内の投票の総数を表します。
- `since`: 集計の開始時刻を表します。

## API エンドポイント

### 場所 (Locations)

#### GET /locations/{group}

すべての場所の一覧を取得します。

- **レスポンス**: Location オブジェクトの配列

#### GET /locations/{group}/{id}

特定の場所の詳細を取得します。

- **レスポンス**: Location オブジェクト

#### POST /locations/{group}

監視対象となる新しい場所を作成します。

- **ボディ**: `id`, `name`, `description`
  - `group`と`id`で場所を一意に識別します。
- **レスポンス**: 作成された Location オブジェクト

### 投票 (Votes)

#### POST /votes/{group}/{id}

新しい混雑状況を報告します。

- **ボディ**: `level` (0-3), `client_address`（ユーザーのIPアドレス）
- **レスポンス**: 作成された Vote オブジェクト

#### GET /votes/{group}/{id}

特定の場所に対する投票結果を取得します。

- **クエリパラメータ**: `before`, `unit`（何分（秒）前から数えるか）
  - クエリパラメータが指定されなければ `{ "before": 10, "unit": "minute" }` として扱います。
- **レスポンス**: `CrowdRate`オブジェクト
  - `group`, `id`, `votes`, `total`, `since`
  - `votes` は float 型で、指定された時間枠内の投票の加重平均を表します。
  - `total` は指定された時間枠内の投票の総数を表します。
  - `since` は集計の開始時刻を表します。

#### GET /votes/{group}

- **クエリパラメータ**: `before`, `unit`（何分（秒）前から数えるか）
  - クエリパラメータが指定されなければ `{ "before": 10, "unit": "minute" }` として扱います。
- **レスポンス**: `CrowdRate`オブジェクトの配列

#### GET /histories/{group}/{id}

特定の場所の投票履歴を取得します。

- **クエリパラメータ**: `before`, `unit`（何分（秒）前から数えるか）
- **レスポンス**: `CrowdRate`オブジェクトの配列

## ビジネスロジック

- 場所の `current_level` は、直近の時間枠（例：過去30分間）に行われた投票の加重平均または最頻値として計算されます。
