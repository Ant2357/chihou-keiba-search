<script>
  import Fa from 'svelte-fa'
  import { faSearch } from '@fortawesome/free-solid-svg-icons'
  import logo from './assets/images/logo-universal.png'
  import {WebScrapingRace} from "../wailsjs/go/main/App.js"

  let url = ""
  let searchResult = {}
  let raceResults = []
  let selectHorse = {
    name: "",
    imgUrl: "",
    pedigree: {}
  }
  let message = ""

  function webScrapingRace() {
      message = "ロード中";
    WebScrapingRace(url).then(result => {
      searchResult = result;
      message = "読み込み完了";
    }).catch(error => {
      message = error;
    })
  }

  function readRaceResults(index) {
    raceResults = searchResult.horses[index].results;
    selectHorse.name = searchResult.horses[index].name;
    selectHorse.imgUrl = searchResult.horses[index].img_url;
    selectHorse.pedigree = searchResult.horses[index].pedigree;
  }
</script>

<div class="container-top container vh-100">
  <div class="row vh-100">
    <div class="col align-self-center">
      <div class="card-rotate card shadow">
        <div class="card-body card-rotate-text">
          <div class="text-center pt-4 mb-3">
            <img
              id="logo"
              alt="Wails logo"
              src="{logo}"
              class="card-img-top h-50 w-50"
            >
          </div>

          <div class="container">
            <div class="text-center">
              <h1 class="display-4">こんにちは!</h1>
              <p class="lead">netkeibaのレースURL情報を打ち込むと、優先的に見たい情報が表示されます</p>
            </div>
          </div>

          <div>
            <div class="text-center mb-3">
              {#if message === "ロード中"}
                <div class="spinner-border" role="status">
                  <span class="visually-hidden">Loading...</span>
                </div>
              {/if}
              <p>{message}</p>
            </div>

            <div class="input-group mb-3">
              <input
                type="text"
                class="form-control"
                placeholder="URLを入力"
                bind:value={url}
              >
              <button
                class="btn btn-outline-success"
                type="button"
                on:click={webScrapingRace}
              >
                <Fa icon={faSearch} />
                <i class="fas fa-search"></i> 生成
              </button>
            </div>

            <div class="mb-3">
              {#if Object.keys(searchResult).length !== 0}
                <div class="mb-3">
                  <h4>{searchResult.name}({searchResult.racetrack}{searchResult.distance})</h4>
                </div>

                <table class="table align-middle">
                  <thead>
                    <tr>
                      <th scope="col">#</th>
                      <th scope="col">名前</th>
                      <th scope="col">レース数</th>
                      <th scope="col">勝利数</th>
                      <th scope="col">敗北数</th>
                      <th scope="col">コース適正</th>
                      <th scope="col">距離適正</th>
                      <th scope="col">脚質</th>
                      <th scope="col">重馬場</th>
                    </tr>
                  </thead>
                  <tbody>
                    {#each searchResult.horses as horse, i}
                      <tr>
                        <th>{i + 1}</th>
                        <td><a href="#race_results" class="link-primary" on:click={() => readRaceResults(i)}>{horse.name}</a></td>
                        <td>{horse.play_game_count}</td>
                        <td>{horse.win}</td>
                        <td>{horse.lose}</td>
                        <td>{horse.course_aptitude}</td>
                        <td>{horse.distance_aptitude}</td>
                        <td>{horse.running_style}</td>
                        <td>{horse.heavy_racetrack}</td>
                      </tr>
                    {/each}
                  </tbody>
                </table>
              {/if}
            </div>

            <div id="race_results" class="mb-3">
              {#if Object.keys(raceResults).length !== 0}
                <div class="mb-3">
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
                        <tr class="{searchResult.distance === `${race.distance}m` ? 'table-success' : ''}">
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
                </div>

                <div class="mb-3 w-100 d-flex justify-content-center">
                  <div class="card w-50 text-center shadow">
                    <div class="card-body">
                      <h5 class="card-title">{selectHorse.name}</h5>

                      <img src={selectHorse.imgUrl} alt="馬の画像" class="horse-img card-img-bottom">

                      <div class="pedigree">
                        <div class="pedigree-item pedigree-parents pedigree-item-father">{selectHorse.pedigree.father}</div>
                        <div class="pedigree-item pedigree-grandparents pedigree-item-father">{selectHorse.pedigree.paternal_grandfather}</div>
                        <div class="pedigree-item pedigree-grandparents pedigree-item-mother">{selectHorse.pedigree.paternal_grandmother}</div>
                        <div class="pedigree-item pedigree-parents pedigree-item-mother">{selectHorse.pedigree.mother}</div>
                        <div class="pedigree-item pedigree-grandparents pedigree-item-father">{selectHorse.pedigree.maternal_grandfather}</div>
                        <div class="pedigree-item pedigree-grandparents pedigree-item-mother">{selectHorse.pedigree.maternal_grandmother}</div>
                      </div>

                    </div>
                  </div>
                </div>
              {/if}
            </div>
          </div>

        </div>
      </div>
    </div>
  </div>
</div>

<style>
.card-rotate {
  transform: rotate(2deg);
}
.card-rotate-text {
  transform:skew(0deg, -2deg);
}

.container-top {
  max-width: 960px;
}

.pedigree {
  display: grid;
  grid-template-rows: repeat(2, auto);
  border: solid 1px #212529;
}
.pedigree-item {
  border: solid 1px #212529;
  display: flex;
  align-items: center;
  justify-content: center;
}
.pedigree-parents {
  grid-column: 1;
  grid-row: span 2;
}
.pedigree-grandparents {
  grid-column: 2;
}
.pedigree-item-father {
  background-color: #a9ceec;
}
.pedigree-item-mother {
  background-color: #f09199;
}

</style>
