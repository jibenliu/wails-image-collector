declare global {
    interface Window {
        backend: Backend;
    }
}

interface Backend {
    MyStruct: {
        Hello():Promise<string>,
        Rename():Promise<string>,
        StoreCount(num:number):Promise<number>,
        RandomValue(num:number):Promise<string>,
        OpenFile():Promise<string>,
        AddUser(name:string):Promise<string>,
        MyBoundMethod(name:string):Promise<string>,
    };
    basic(): Promise<string>;

    NetWorkStatus(): Promise<boolean>;

    StoreFrontend(num:number) :Promise<boolean>
}

export default global
