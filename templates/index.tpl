{{define "body"}}
<div class="jumbotron">
  <h1>SCVL URL Shortener</h1>
  {{if .LoginURL}}
    <a href="{{.LoginURL}}">ログイン</a>
  {{end}}
  <p>URLの短縮ができます。</p>
  <form action="/shorten" method="post" class="form-inline">
    <input type="url" name="url" value="{{.URL}}" class="form-control">
    <input type="submit" value="送信" class="btn btn-primary">
  </form>
  {{if .URL}}
    <div class="panel panel-success">
      <div class="panel-heading">
        <h3 class="panel-title">短縮が成功しました。</h3>
      </div>
      <div class="panel-body">
        <p>短縮結果:</p>
        <p><script>document.write(location.origin)</script>/{{.Slug}}</p>
      </div>
    </div>
  {{end}}
</div>
{{if .User}}
  <h2>{{.User.Name}}の短縮したURL一覧</h2>
  <table class="table urls">
    <tr>
      <th width="80">短縮URL</th>
      <th>リダイレクト先URL</th>
      <th width="100">クリック数</th>
    </tr>
    {{range .User.Pages}}
      <tr>
        <td><a href="/{{.Slug}}" target="_blank">/{{.Slug}}</a></td>
        <td class="truncate"><a href="{{.URL}}" target="_blank">{{.URL}}</a></td>
        <td>{{.ViewCount}}</td>
      </tr>
    {{end}}
  </table>
{{end}}
{{end}}
