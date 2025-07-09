export namespace stock {
	
	export class StockData {
	    name: string;
	    code: string;
	    price: number;
	    change: number;
	    changePct: number;
	    volume: number;
	    amount: number;
	    open: number;
	    prevClose: number;
	    high: number;
	    low: number;
	    bid: number;
	    ask: number;
	    bidVolume: number;
	    askVolume: number;
	    marketStatus: string;
	    // Go type: time
	    updateTime: any;
	    date: string;
	    time: string;
	
	    static createFrom(source: any = {}) {
	        return new StockData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.code = source["code"];
	        this.price = source["price"];
	        this.change = source["change"];
	        this.changePct = source["changePct"];
	        this.volume = source["volume"];
	        this.amount = source["amount"];
	        this.open = source["open"];
	        this.prevClose = source["prevClose"];
	        this.high = source["high"];
	        this.low = source["low"];
	        this.bid = source["bid"];
	        this.ask = source["ask"];
	        this.bidVolume = source["bidVolume"];
	        this.askVolume = source["askVolume"];
	        this.marketStatus = source["marketStatus"];
	        this.updateTime = this.convertValues(source["updateTime"], null);
	        this.date = source["date"];
	        this.time = source["time"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

