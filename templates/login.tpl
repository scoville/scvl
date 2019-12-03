{{define "login"}}
<div class="login">
  <p>ログインしてください。</p>
  <form action="/login" method="post">
    <div class="form-inline">
      <input id="email" type="email" name="email" class="form-control" placeholder="メールアドレス">
      <input id="password" type="password" name="password" class="form-control" placeholder="パスワード">
      <input type="submit" value="ログイン" class="btn btn-primary">
    </div>
  </form>
  <p>または、</p>
  <a href="{{.LoginURL}}" class="login btn btn-primary btn-lg">Googleアカウントでログイン</a>
</div>
{{end}}