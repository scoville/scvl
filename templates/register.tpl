{{define "body"}}
<div class="jumbotron">
  <h1>ユーザー登録をする。</h1>
  <form action="/register" method="post">
    <div class="form">
      <p>メールアドレス： {{.Email}}</p>
      <p>パスワード：<input id="password" type="password" name="password" class="form-control"></p>
      <input type="hidden" name="hash" value="{{.Hash}}">
      <p><input type="submit" value="送信" class="btn btn-primary"></p>
    </div>
  </form>
</div>
{{end}}
