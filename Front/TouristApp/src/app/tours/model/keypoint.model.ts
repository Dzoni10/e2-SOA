export interface Keypoint {
    id? : string;
    tourId? : string;  
  name: string;
  description: string;
  latitude: number;
  longitude: number;
  images: string[];
  formData?: FormData;
}