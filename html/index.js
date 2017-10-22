var root=new Vue({
    el: "#root",
    data: {
        ids:[],
        videos:[],
    },
    methods:{
      change(e){
          this.videos=[]
          this.ids=$(e.target).val().split("\n").map((u)=>u.substr(u.indexOf("v=")+2)).filter((u)=>u.length>0)
          for(let id of this.ids){
              $.get("video/"+id,function (res) {
                  res.select=0
                  root.videos.push(res)
              })
          }
      },
      down(){
        console.log(this.videos)
        for(let v of this.videos){
            let f=v.Formats[v.select]
            let url=`${location.href}dl/?fname=${encodeURIComponent(v.Title)}_${f.Res}.${f.Ext}&clen=${f.Clen}&url=${btoa(f.Url)}`
            window.open(url)
            console.log(url)
        }
      }
    },
    computed:{
        count:function () {
            let c=0
            for(let v of this.videos){
                c+=parseInt(v.Formats[v.select].Clen)
            }
            return parseInt(c/1024/1024)
        }
    }
})
