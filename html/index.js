var root=new Vue({
    el: "#root",
    data: {
        ids:[],
        videos:[],
    },
    methods:{
      change:function(e){
          var urls=$(e.target).val()
          var k="watch?v="
          this.ids=urls
              .split("\n")
              .map( function (u) {return u.substring(u.indexOf(k)+k.length)})
              .filter(function (u) {return u.length>0})
          console.log(urls)
          this.videos=[]
          for(var i in this.ids){
              $.get("video/"+this.ids[i],function (res) {
                  if(res["Formats"][0].Res=="1080p"){
                      res.select=1
                  }else{
                      res.select=0
                  }
                  root.videos.push(res)
              })
          }
      },
      down:function(){
        console.log(this.videos)
        for(var i in this.videos){
            var v=this.videos[i]
            var f=v.Formats[v.select]
            var url=location.href+"video/"+v.ID+"/format/"+f.Itag
            window.open(url)
        }
      }
    },
    computed:{
        count:function () {
            var c=0
            for(var i in this.videos){
                var v=this.videos[i]
                c+=parseInt(v.Formats[v.select].Clen)
            }
            return parseInt(c/1024/1024)
        }
    }
})
