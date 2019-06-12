{{define "body"}}
<div class="jumbotron">
  <h1>SCVL Image Uploader</h1>
  {{if .LoginURL}}
    <a href="{{.LoginURL}}" class="login btn btn-primary btn-lg">ログイン</a>
  {{else}}
    <p>画像のアップロードができます。</p>
    <form action="/images" method="post" enctype="multipart/form-data">
      <input type="file" name="file" required>
      <p style="font-size: 12px;">※ ファイルサイズの上限は50MBです。</p>
      <div class="mt20">
        <input type="submit" value="送信" class="btn btn-primary">
      </div>
    </form>
  {{end}}
  {{if .URL}}
    <div class="panel panel-success">
      <div class="panel-heading">
        <h3 class="panel-title">画像のアップロードが成功しました。</h3>
      </div>
      <div class="panel-body">
        <p>画像表示用リンク:</p>
        <p>
          <input id="urlImage" class="copy-target" type="text" value="{{.URL}}" readonly>
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
  <h2>{{.User.Name}}のアップロードした画像一覧</h2>
  <table class="table files">
    <tr>
      <th width="80">URL</th>
    </tr>
    {{range .User.Images}}
      <tr>
        <td class="truncate"><a href="{{.URL}}" target="_blank">{{.URL}}</a></td>
      </tr>
    {{end}}
  </table>
{{end}}
{{end}}
