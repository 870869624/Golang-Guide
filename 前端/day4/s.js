function example() {
  const name = "张三";  // 普通变量 - "忠诚的"
  console.log(name);    // 永远是"张三"
  
  console.log(this);    // this - "善变的"
  // 这个this是谁，取决于函数怎么被调用
}