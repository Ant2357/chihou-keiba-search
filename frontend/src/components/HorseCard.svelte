<script>
  export let horse;
</script>

<div class="card w-50 text-center shadow">
  <div class="card-body">
    <h5 class="card-title">{horse.name}</h5>

    <img src={horse.imgUrl} alt="馬の画像" class="horse-img card-img-bottom">

    <div class="pedigree">
      <div class="pedigree-item pedigree-parents pedigree-item-father">{horse.pedigree.father}</div>
      <div class="pedigree-item pedigree-grandparents pedigree-item-father">{horse.pedigree.paternal_grandfather}</div>
      <div class="pedigree-item pedigree-grandparents pedigree-item-mother">{horse.pedigree.paternal_grandmother}</div>
      <div class="pedigree-item pedigree-parents pedigree-item-mother">{horse.pedigree.mother}</div>
      <div class="pedigree-item pedigree-grandparents pedigree-item-father">{horse.pedigree.maternal_grandfather}</div>
      <div class="pedigree-item pedigree-grandparents pedigree-item-mother">{horse.pedigree.maternal_grandmother}</div>
    </div>

  </div>
</div>

<style>
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
