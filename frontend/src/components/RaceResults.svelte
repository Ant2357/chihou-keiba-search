<script>
  export let raceResults; 
  export let activeDistance; 
</script>

<table class="table">
  <thead>
    <tr>
      <th scope="col">日付</th>
      <th scope="col">レース名</th>
      <th scope="col">着順</th>
      <th scope="col">距離</th>
      <th scope="col">馬場</th>
      <th scope="col">タイム</th>
    </tr>
  </thead>
  <tbody>
    {#each raceResults as race}
      <tr class="{activeDistance === `${race.distance}m` ? 'table-success' : ''}">
        <td>{race.date.replaceAll(/T.*/g, "").replaceAll("-", "/")}</td>
        <td>{race.raceName}</td>
        <td>{race.result === -1 ? "除" : race.result}</td>
        <td>{race.distance}</td>
        <td>{race.baba}</td>
        <td>{race.time}</td>
      </tr>
    {/each}
  </tbody>
</table>
