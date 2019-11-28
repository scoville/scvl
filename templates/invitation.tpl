{{define "body"}}
<div class="jumbotron">
  <div class="panel panel-success">
    <div class="panel-heading">
      <p>ユーザーを招待しました。以下のURLを共有し、24時間以内にユーザー登録を完了してください。</p>
    </div>
    <div class="panel-body">
      <p>登録URL:</p>
      <p>
        <input id="shortenUrl" class="copy-target" type="text" value="{{.RegisterPath}}" readonly>
        <span class="copy">
          <i class="material-icons">content_copy</i>
          <span class="hint">コピーする</span>
        </span>
      </p>
    </div>
  </div>
  <ul><li><a href="{{.MainHost}}">戻る</a></li></ul>
</div>
{{end}}
