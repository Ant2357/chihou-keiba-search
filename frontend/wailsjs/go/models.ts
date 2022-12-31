export namespace web_scraping {
	
	export class RaceResult {
	    // Go type: time.Time
	    date: any;
	    raceName: string;
	    result: number;
	    distance: string;
	    baba: string;
	    time: string;
	
	    static createFrom(source: any = {}) {
	        return new RaceResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = this.convertValues(source["date"], null);
	        this.raceName = source["raceName"];
	        this.result = source["result"];
	        this.distance = source["distance"];
	        this.baba = source["baba"];
	        this.time = source["time"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class Horse {
	    img_url: string;
	    name: string;
	    play_game_count: number;
	    win: number;
	    lose: number;
	    course_aptitude: string;
	    distance_aptitude: string;
	    running_style: string;
	    heavy_racetrack: string;
	    results: RaceResult[];
	
	    static createFrom(source: any = {}) {
	        return new Horse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.img_url = source["img_url"];
	        this.name = source["name"];
	        this.play_game_count = source["play_game_count"];
	        this.win = source["win"];
	        this.lose = source["lose"];
	        this.course_aptitude = source["course_aptitude"];
	        this.distance_aptitude = source["distance_aptitude"];
	        this.running_style = source["running_style"];
	        this.heavy_racetrack = source["heavy_racetrack"];
	        this.results = this.convertValues(source["results"], RaceResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class Race {
	    name: string;
	    racetrack: string;
	    distance: string;
	    horses: Horse[];
	
	    static createFrom(source: any = {}) {
	        return new Race(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.racetrack = source["racetrack"];
	        this.distance = source["distance"];
	        this.horses = this.convertValues(source["horses"], Horse);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

