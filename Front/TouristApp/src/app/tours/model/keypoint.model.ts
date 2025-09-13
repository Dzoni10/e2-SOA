export interface Keypoint {
    id? : string;
    tourId? : string;  
  name: string;
  description: string;
  latitude: number;
  longitude: number;
  order?: number;
  images: string[];
  formData?: FormData;
  marker?: L.Marker; // optional property
}