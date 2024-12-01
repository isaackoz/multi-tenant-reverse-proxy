export type TApiData<T> =
	| {
			data: T;
			error?: null;
	  }
	| {
			data?: null;
			error: {
				message: string;
			};
	  };
