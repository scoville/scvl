{{define "body"}}
<div class="email-preview">
  <h1>メールプレビュー</h1>

  <p>以下のメールを{{.SendingNumber}}に送信します。（最初の３件のみ表示）</p>

  <table class="table">
    <tr>
      <th>送信先</th>
      <td>送信文言</td>
    </tr>
    {{range .Emails}}
      <tr>
        <td></td>
        <td></td>
      </tr>
    {{end}}
  </table>
  <form action="{{if .IsEmailDomain}}/{{else}}/emails{{end}}" method="post">
    <input type="hidden" name="spreadsheet_url" value="{{.SpreadsheetURL}}" />
    <input type="hidden" name="sheet_name" value="{{.SheetName}}" />
    <input type="hidden" name="sender" value="{{.Sender}}" />
    <input type="hidden" name="template" value="{{.Template}}" />
    <input type="submit" value="送信" class="btn btn-primary">
  </form>
</div>
{{end}}
