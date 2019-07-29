{{define "body"}}
<div class="jumbotron">
  <h1>SCVL URL Shortener</h1>
  {{if .LoginURL}}
    <a href="{{.LoginURL}}" class="login btn btn-primary btn-lg">ログイン</a>
  {{else}}
    <p>URLの短縮ができます。</p>
    <form action="/shorten" method="post">
      <div class="form-inline">
        <input id="url" type="url" name="url" value="{{.URL}}" class="form-control">
        <input type="submit" value="送信" class="btn btn-primary">
      </div>
      <label class="toggle-utm">
        <input id="utm" type="checkbox" name="utm">トラッキング用パラメータの追加
      </label>
      <div class="utm-information">
        <div class="form-group">
          <label for="utm_source">【必須】utm_source（媒体名）</label>
          <input class="form-control" type="text" name="utm_source" value="" placeholder="google, facebook, line など">
        </div>
        <div class="form-group">
          <label for="utm_medium">【必須】utm_medium（メディアの種類）</label>
          <select class="form-control" name="utm_medium">
            <option value="-"></option>
            <option value="paidsearch">paidsearch（有料検索）</option>
            <option value="display">display（ディスプレイ広告）</option>
            <option value="affiliate">affiliate（その他の広告）</option>
            <option value="social">social（SNS）</option>
            <option value="email">email（メール）</option>
          </select>
        </div>
        <div class="form-group">
          <label for="utm_campaign">【任意】utm_campaign(キャンペーン名)</label>
          <input class="form-control" type="text" name="utm_campaign" value="" placeholder="spring_sale など">
        </div>
      </div>

      <label class="toggle-ogp">
        <input id="ogp" type="checkbox" name="ogp">OGPをカスタマイズ
      </label>
      <div class="ogp-information">
        <div class="form-group">
          <label for="title">タイトル</label>
          <input class="ogp form-control" type="text" name="title" value="">
        </div>
        <div class="form-group">
          <label for="title">画像URL</label>
          <input class="ogp form-control" type="url" name="image" value="">
        </div>
        <div class="form-group">
          <label for="title">説明文</label>
          <input class="ogp form-control" type="text" name="description" value="">
        </div>
      </div>
    </form>
  {{end}}

  {{if .URL}}
    <div class="panel panel-success">
      <div class="panel-heading">
        <h3 class="panel-title">短縮が成功しました。</h3>
      </div>
      <div class="panel-body">
        <p>短縮結果:</p>
        <p>
          <input id="shortenUrl" class="copy-target" type="text" value="{{.Slug}}" readonly>
          <span class="copy">
            <i class="material-icons">content_copy</i>
            <span class="hint">コピーする</span>
          </span>
        </p>
      </div>
    </div>
  {{end}}
</div>

{{if .User}}
  <h2>{{.User.Name}}の短縮したURL一覧</h2>
  <p>全{{.Count}}件中、{{.Start}}〜{{.End}}を表示中</p>
  <table class="table urls">
    <tr>
      <th width="80">短縮URL</th>
      <th>リダイレクト先URL</th>
      <th width="100">QRコード</th>
      <th width="100">クリック数</th>
      <th width="100">編集</th>
    </tr>
    {{range .Pages}}
      <tr>
        <td><a href="/{{.Slug}}" target="_blank">/{{.Slug}}</a></td>
        <td>
          <div class="truncate">
            <a href="{{.URL}}" target="_blank">
              {{.URL}}
            </a>
          </div>
          <div class="truncate mt10">{{.Title}}</div>
        </td>
        <td class="qr">
          <a href="/{{.Slug}}/qr.png" target="_blank">
            <img src="/{{.Slug}}/qr.png">
          </a>
        </td>
        <td>{{.ViewCount}}</td>
        <td>
          <a href="/{{.Slug}}/edit" class="btn btn-default">編集</a>
        </td>
      </tr>
    {{end}}
  </table>
  {{if .PrevURL}}
    <a href="{{.PrevURL}}" class="btn btn-default">前のページへ</a>
  {{end}}
  {{if .NextURL}}
    <a href="{{.NextURL}}" class="btn btn-default">次のページへ</a>
  {{end}}
{{end}}
{{end}}
