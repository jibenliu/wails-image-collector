
async function StoreFrontend(num: number) {
    window.backend.MyStruct.StoreCount(num).then((result:number) => {
        console.log(result)
    });
}

export default StoreFrontend
