{{define "body"}}
<div class="jumbotron">
  <h1>SCVL File Uploader</h1>
  {{if .LoginURL}}
    <a href="{{.LoginURL}}" class="login btn btn-primary btn-lg">ログイン</a>
  {{else}}
    <p>ファイルのアップロードができます。</p>
    <form action="/files" method="post" enctype="multipart/form-data">
      <div class="form-group">
        <label for="download_limit">ダウンロード制限回数（無制限の場合は0を指定）</label>
        <input type="number" name="download_limit" value="1" class="form-control" required />
      </div>
      <div class="form-group">
        <label for="valid_days">ダウンロード可能日数（無期限の場合は0を指定）</label>
        <input type="number" name="valid_days" value="0" class="form-control" required />
      </div>
      <div class="form-group">
        <label for="email">送信先メールアドレス（任意）</label>
        <input type="email" name="email" value="" class="form-control" reuired />
      </div>
      <div class="form-group">
        <label for="password">パスワード（任意）</label>
        <input type="password" name="password" value="" class="form-control" />
      </div>
      <input type="file" name="file" required>
      <div class="mt20">
        <input type="submit" value="送信" class="btn btn-primary">
      </div>
    </form>
  {{end}}
  {{if .Slug}}
    <div class="panel panel-success">
      <div class="panel-heading">
        <h3 class="panel-title">ファイルのアップロードが成功しました。</h3>
      </div>
      <div class="panel-body">
        <p>ダウンロード用リンク:</p>
        <p>
          <input id="shortenUrl" type="text" value="{{.Slug}}" readonly>
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
  <h2>{{.User.Name}}のアップロードしたファイル一覧</h2>
  <table class="table files">
    <tr>
      <th width="80">ダウンロード用リンク</th>
      <th width="140">ダウンロード回数 / 制限回数</th>
      <th width="100">ダウンロード期限</th>
      <th width="100">編集</th>
    </tr>
    {{range .User.Files}}
      <tr>
        <td class="truncate"><a href="/files/{{.Slug}}" target="_blank">/files/{{.Slug}}</a></td>
        <td>{{.DownloadCount}}{{if ne .DownloadLimit 0}} / {{.DownloadLimit}}{{end}}</td>
        <td>{{.FormatDeadline}}</td>
        <td>
          <a href="/files/{{.Slug}}/edit" class="btn btn-default">編集</a>
        </td>
      </tr>
    {{end}}
  </table>
{{end}}
{{end}}
